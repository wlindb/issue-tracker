package keycloak_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/wlindb/issue-tracker/internal/pkg/keycloak"
)

func localJWKSServer(t *testing.T, publicKey *rsa.PublicKey) *httptest.Server {
	t.Helper()
	n := base64.RawURLEncoding.EncodeToString(publicKey.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(publicKey.E)).Bytes())
	body, err := json.Marshal(map[string]any{
		"keys": []map[string]any{
			{"kty": "RSA", "alg": "RS256", "use": "sig", "n": n, "e": e},
		},
	})
	require.NoError(t, err)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}))
	t.Cleanup(server.Close)
	return server
}

func mintToken(t *testing.T, privateKey *rsa.PrivateKey, claims jwt.Claims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(privateKey)
	require.NoError(t, err)
	return signed
}

func Test_KeycloakTokenValidator_ValidToken_ReturnsUserID(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	jwksServer := localJWKSServer(t, &privateKey.PublicKey)

	validator, err := keycloak.NewKeycloakTokenValidator(jwksServer.URL)
	require.NoError(t, err)

	expected := uuid.New()
	rawToken := mintToken(t, privateKey, jwt.RegisteredClaims{
		Subject:   expected.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	})

	actual, err := validator.ValidateToken(context.Background(), rawToken)

	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func Test_KeycloakTokenValidator_ExpiredToken_ReturnsError(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	jwksServer := localJWKSServer(t, &privateKey.PublicKey)

	validator, err := keycloak.NewKeycloakTokenValidator(jwksServer.URL)
	require.NoError(t, err)

	rawToken := mintToken(t, privateKey, jwt.RegisteredClaims{
		Subject:   uuid.New().String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
	})

	_, err = validator.ValidateToken(context.Background(), rawToken)

	require.Error(t, err)
}

func Test_KeycloakTokenValidator_WrongKey_ReturnsError(t *testing.T) {
	jwksKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	signingKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	jwksServer := localJWKSServer(t, &jwksKey.PublicKey)

	validator, err := keycloak.NewKeycloakTokenValidator(jwksServer.URL)
	require.NoError(t, err)

	rawToken := mintToken(t, signingKey, jwt.RegisteredClaims{
		Subject:   uuid.New().String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	})

	_, err = validator.ValidateToken(context.Background(), rawToken)

	require.Error(t, err)
}

func Test_KeycloakTokenValidator_MalformedToken_ReturnsError(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	jwksServer := localJWKSServer(t, &privateKey.PublicKey)

	validator, err := keycloak.NewKeycloakTokenValidator(jwksServer.URL)
	require.NoError(t, err)

	_, err = validator.ValidateToken(context.Background(), "not-a-jwt")

	require.Error(t, err)
}

func Test_KeycloakTokenValidator_NonUUIDSub_ReturnsError(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	jwksServer := localJWKSServer(t, &privateKey.PublicKey)

	validator, err := keycloak.NewKeycloakTokenValidator(jwksServer.URL)
	require.NoError(t, err)

	rawToken := mintToken(t, privateKey, jwt.RegisteredClaims{
		Subject:   "not-a-uuid",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	})

	_, err = validator.ValidateToken(context.Background(), rawToken)

	require.Error(t, err)
}

func Test_KeycloakTokenValidator_MissingSub_ReturnsError(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	jwksServer := localJWKSServer(t, &privateKey.PublicKey)

	validator, err := keycloak.NewKeycloakTokenValidator(jwksServer.URL)
	require.NoError(t, err)

	rawToken := mintToken(t, privateKey, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	_, err = validator.ValidateToken(context.Background(), rawToken)

	require.Error(t, err)
}
