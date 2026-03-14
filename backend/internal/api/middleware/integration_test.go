//go:build integration

package middleware_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/wlindb/issue-tracker/internal/api"
	"github.com/wlindb/issue-tracker/internal/api/middleware"
)

type keycloakContainer struct {
	container testcontainers.Container
	baseURL   string
}

var kc *keycloakContainer

func TestMain(m *testing.M) {
	ctx := context.Background()
	var err error
	kc, err = startKeycloakContainer(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to start keycloak container: %v\n", err)
		os.Exit(1)
	}
	defer kc.container.Terminate(ctx) //nolint:errcheck
	os.Exit(m.Run())
}

func startKeycloakContainer(ctx context.Context) (*keycloakContainer, error) {
	req := testcontainers.ContainerRequest{
		Image: "quay.io/keycloak/keycloak:26",
		Cmd:   []string{"start-dev", "--import-realm"},
		Env: map[string]string{
			"KEYCLOAK_ADMIN":          "admin",
			"KEYCLOAK_ADMIN_PASSWORD": "admin",
		},
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      "../../keycloak/realm-export.json",
				ContainerFilePath: "/opt/keycloak/data/import/realm-export.json",
				FileMode:          0o644,
			},
		},
		ExposedPorts: []string{"8080/tcp"},
		WaitingFor: wait.ForHTTP("/realms/issue-tracker").
			WithPort("8080").
			WithStartupTimeout(120 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("start container: %w", err)
	}

	mappedPort, err := container.MappedPort(ctx, "8080")
	if err != nil {
		return nil, fmt.Errorf("mapped port: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("host: %w", err)
	}

	baseURL := fmt.Sprintf("http://%s:%s", host, mappedPort.Port())
	return &keycloakContainer{container: container, baseURL: baseURL}, nil
}

func keycloakJWKSURL() string {
	return kc.baseURL + "/realms/issue-tracker/protocol/openid-connect/certs"
}

// obtainToken fetches a real JWT from Keycloak via Resource Owner Password grant.
func obtainToken(t *testing.T) string {
	t.Helper()
	tokenURL := kc.baseURL + "/realms/issue-tracker/protocol/openid-connect/token"
	data := url.Values{
		"grant_type":    {"password"},
		"client_id":     {"issue-tracker-app"},
		"client_secret": {"test-secret"},
		"username":      {"testuser"},
		"password":      {"password"},
	}
	resp, err := http.PostForm(tokenURL, data)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "token request failed: %s", body)

	var result map[string]any
	require.NoError(t, json.Unmarshal(body, &result))
	token, ok := result["access_token"].(string)
	require.True(t, ok, "access_token missing from response")
	return token
}

func newProtectedEcho(jwksURL string) *echo.Echo {
	e := echo.New()
	e.Use(middleware.JwtMiddleware(jwksURL))
	e.Use(middleware.UserIDMiddleware())
	e.GET("/ping", func(c echo.Context) error {
		id := api.UserIDFromContext(c.Request().Context())
		return c.String(http.StatusOK, id.String())
	})
	return e
}

// localJWKSServer creates an httptest.Server serving the public key as a JWKS.
func localJWKSServer(t *testing.T, pub *rsa.PublicKey) *httptest.Server {
	t.Helper()
	n := base64.RawURLEncoding.EncodeToString(pub.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(pub.E)).Bytes())
	body, err := json.Marshal(map[string]any{
		"keys": []map[string]any{
			{"kty": "RSA", "alg": "RS256", "use": "sig", "n": n, "e": e},
		},
	})
	require.NoError(t, err)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(body) //nolint:errcheck
	}))
	t.Cleanup(srv.Close)
	return srv
}

func Test_Integration_ValidKeycloakToken_PassesMiddleware(t *testing.T) {
	token := obtainToken(t)

	e := newProtectedEcho(keycloakJWKSURL())
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	sub := strings.TrimSpace(rec.Body.String())
	require.Regexp(t, `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, sub)
}

func Test_Integration_NoToken_Returns401(t *testing.T) {
	e := newProtectedEcho(keycloakJWKSURL())
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func Test_Integration_WrongIssuerToken_Returns401(t *testing.T) {
	// Token signed by a key not in Keycloak's JWKS — signature validation fails.
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
		Subject:   "00000000-0000-0000-0000-000000000002",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	})
	token, err := tok.SignedString(key)
	require.NoError(t, err)

	e := newProtectedEcho(keycloakJWKSURL())
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func Test_Integration_ExpiredToken_Returns401(t *testing.T) {
	// Use a local JWKS (not Keycloak's) so we control the key and can set expiry.
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	srv := localJWKSServer(t, &key.PublicKey)

	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
		Subject:   "00000000-0000-0000-0000-000000000001",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
	})
	token, err := tok.SignedString(key)
	require.NoError(t, err)

	e := newProtectedEcho(srv.URL)
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}
