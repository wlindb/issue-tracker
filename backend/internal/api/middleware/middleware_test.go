//go:build !integration

package middleware_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/wlindb/issue-tracker/internal/api"
	"github.com/wlindb/issue-tracker/internal/api/middleware"
)

func generateTestKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	return key
}

func serveJWKS(t *testing.T, pub *rsa.PublicKey) *httptest.Server {
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
		_, err = w.Write(body)
		require.NoError(t, err)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func mintToken(t *testing.T, priv *rsa.PrivateKey, claims jwt.RegisteredClaims) string {
	t.Helper()
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := tok.SignedString(priv)
	require.NoError(t, err)
	return signed
}

func newPingServer(t *testing.T, jwksURL string) *echo.Echo {
	t.Helper()
	e := echo.New()
	e.Use(middleware.JwtMiddleware(jwksURL))
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	return e
}

func validClaims() jwt.RegisteredClaims {
	return jwt.RegisteredClaims{
		Subject:   "00000000-0000-0000-0000-000000000001",
		Issuer:    "issue-tracker",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}
}

func Test_JwtMiddleware_NoHeader(t *testing.T) {
	key := generateTestKey(t)
	srv := serveJWKS(t, &key.PublicKey)

	e := newPingServer(t, srv.URL)
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func Test_JwtMiddleware_BadScheme(t *testing.T) {
	key := generateTestKey(t)
	srv := serveJWKS(t, &key.PublicKey)

	e := newPingServer(t, srv.URL)
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Authorization", "Basic abc")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func Test_JwtMiddleware_InvalidToken(t *testing.T) {
	key := generateTestKey(t)
	srv := serveJWKS(t, &key.PublicKey)

	e := newPingServer(t, srv.URL)
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Authorization", "Bearer garbage")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func Test_JwtMiddleware_ExpiredToken(t *testing.T) {
	key := generateTestKey(t)
	srv := serveJWKS(t, &key.PublicKey)

	claims := validClaims()
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-1 * time.Hour))
	token := mintToken(t, key, claims)

	e := newPingServer(t, srv.URL)
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func Test_JwtMiddleware_WrongKey(t *testing.T) {
	key := generateTestKey(t)
	srv := serveJWKS(t, &key.PublicKey)

	otherKey := generateTestKey(t)
	token := mintToken(t, otherKey, validClaims())

	e := newPingServer(t, srv.URL)
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func Test_JwtMiddleware_ValidToken(t *testing.T) {
	key := generateTestKey(t)
	srv := serveJWKS(t, &key.PublicKey)

	token := mintToken(t, key, validClaims())

	e := newPingServer(t, srv.URL)
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func newPingServerWithUserID(t *testing.T, jwksURL string) *echo.Echo {
	t.Helper()
	e := echo.New()
	e.Use(middleware.JwtMiddleware(jwksURL))
	e.Use(middleware.UserIDMiddleware())
	e.GET("/ping", func(c echo.Context) error {
		id, err := api.UserIDFromContext(c.Request().Context())
		if err != nil {
			return c.String(http.StatusUnauthorized, err.Error())
		}
		return c.String(http.StatusOK, id.String())
	})
	return e
}

func Test_UserIDMiddleware_PopulatesContext(t *testing.T) {
	key := generateTestKey(t)
	srv := serveJWKS(t, &key.PublicKey)

	token := mintToken(t, key, validClaims())

	e := newPingServerWithUserID(t, srv.URL)
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "00000000-0000-0000-0000-000000000001", rec.Body.String())
}

// — WorkspaceMembershipMiddleware tests —

type mockWorkspaceMemberChecker struct {
	mock.Mock
}

func (m *mockWorkspaceMemberChecker) IsMember(ctx context.Context, workspaceID uuid.UUID, userID uuid.UUID) (bool, error) {
	args := m.Called(ctx, workspaceID, userID)
	return args.Bool(0), args.Error(1)
}

func injectUserID(id uuid.UUID) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := api.WithUserID(c.Request().Context(), id)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func newWorkspaceMiddlewareServer(checker middleware.WorkspaceMemberChecker, userID uuid.UUID) *echo.Echo {
	e := echo.New()
	e.Use(injectUserID(userID))
	e.Use(middleware.WorkspaceMembershipMiddleware(checker))
	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
	return e
}

func Test_WorkspaceMembershipMiddleware_MissingHeader_Returns400(t *testing.T) {
	checker := &mockWorkspaceMemberChecker{}
	e := newWorkspaceMiddlewareServer(checker, uuid.New())

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	checker.AssertNotCalled(t, "IsMember")
}

func Test_WorkspaceMembershipMiddleware_InvalidUUID_Returns400(t *testing.T) {
	checker := &mockWorkspaceMemberChecker{}
	e := newWorkspaceMiddlewareServer(checker, uuid.New())

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Workspace-ID", "not-a-uuid")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	checker.AssertNotCalled(t, "IsMember")
}

func Test_WorkspaceMembershipMiddleware_NonMember_Returns403(t *testing.T) {
	checker := &mockWorkspaceMemberChecker{}
	workspaceID, userID := uuid.New(), uuid.New()
	checker.On("IsMember", mock.Anything, workspaceID, userID).Return(false, nil)

	e := newWorkspaceMiddlewareServer(checker, userID)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Workspace-ID", workspaceID.String())
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusForbidden, rec.Code)
	checker.AssertExpectations(t)
}

func Test_WorkspaceMembershipMiddleware_IsMemberError_Returns500(t *testing.T) {
	checker := &mockWorkspaceMemberChecker{}
	workspaceID, userID := uuid.New(), uuid.New()
	checker.On("IsMember", mock.Anything, workspaceID, userID).Return(false, errors.New("db down"))

	e := newWorkspaceMiddlewareServer(checker, userID)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Workspace-ID", workspaceID.String())
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	checker.AssertExpectations(t)
}

func Test_WorkspaceMembershipMiddleware_ValidMember_CallsNextAndSetsContext(t *testing.T) {
	checker := &mockWorkspaceMemberChecker{}
	workspaceID, userID := uuid.New(), uuid.New()
	checker.On("IsMember", mock.Anything, workspaceID, userID).Return(true, nil)

	var contextWorkspaceID uuid.UUID
	e := echo.New()
	e.Use(injectUserID(userID))
	e.Use(middleware.WorkspaceMembershipMiddleware(checker))
	e.GET("/test", func(c echo.Context) error {
		contextWorkspaceID, _ = api.WorkspaceIDFromContext(c.Request().Context())
		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Workspace-ID", workspaceID.String())
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, workspaceID, contextWorkspaceID)
	checker.AssertExpectations(t)
}

func Test_UserIDMiddleware_NonUUIDSub(t *testing.T) {
	key := generateTestKey(t)
	srv := serveJWKS(t, &key.PublicKey)

	claims := validClaims()
	claims.Subject = "not-a-uuid"
	token := mintToken(t, key, claims)

	e := newPingServerWithUserID(t, srv.URL)
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}
