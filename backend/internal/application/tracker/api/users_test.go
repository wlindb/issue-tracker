//go:build !integration

package api_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	rootapi "github.com/wlindb/issue-tracker/internal/application/api"
	"github.com/wlindb/issue-tracker/internal/application/api/model"
	"github.com/wlindb/issue-tracker/internal/application/tracker/api"
	userdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/user"
)

type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) Upsert(ctx context.Context, command userdomain.UpsertUserCommand) (userdomain.User, error) {
	args := m.Called(ctx, command)
	return args.Get(0).(userdomain.User), args.Error(1)
}

func newUserTestServer(t *testing.T, service api.UserService) *echo.Echo {
	t.Helper()
	e := echo.New()
	h := &api.Handler{
		UserHandler: api.NewUserHandler(service),
	}
	strict := model.NewStrictHandler(&rootapi.Handler{Handler: *h}, nil)
	model.RegisterHandlersWithBaseURL(e, strict, "/api/v1")
	return e
}

func Test_UpsertCurrentUser_ValidClaims_Returns200(t *testing.T) {
	service := &mockUserService{}
	userID := uuid.New()

	service.On("Upsert", mock.Anything, userdomain.UpsertUserCommand{
		ID:    userID,
		Email: "jane@example.com",
		Name:  "Jane Doe",
	}).Return(userdomain.User{
		ID:    userID,
		Email: "jane@example.com",
		Name:  "Jane Doe",
	}, nil)

	e := newUserTestServer(t, service)
	e.Use(injectUser(userID))
	e.Use(injectUserClaims(api.UserClaims{Email: "jane@example.com", Name: "Jane Doe"}))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/me", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got model.User
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	assert.Equal(t, userID, got.Id)
	assert.Equal(t, "Jane Doe", got.Name)
	service.AssertExpectations(t)
}

func Test_UpsertCurrentUser_MissingUserID_Returns401(t *testing.T) {
	service := &mockUserService{}

	e := newUserTestServer(t, service) // no injectUser

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/me", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
	service.AssertNotCalled(t, "Upsert")
}

func Test_UpsertCurrentUser_MissingClaims_Returns401(t *testing.T) {
	service := &mockUserService{}

	e := newUserTestServer(t, service)
	e.Use(injectUser(uuid.New())) // no injectUserClaims

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/me", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
	service.AssertNotCalled(t, "Upsert")
}

func Test_UpsertCurrentUser_ServiceError_Returns500(t *testing.T) {
	service := &mockUserService{}
	userID := uuid.New()

	service.On("Upsert", mock.Anything, mock.AnythingOfType("user.UpsertUserCommand")).
		Return(userdomain.User{}, errors.New("db down"))

	e := newUserTestServer(t, service)
	e.Use(injectUser(userID))
	e.Use(injectUserClaims(api.UserClaims{Email: "jane@example.com", Name: "Jane Doe"}))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/me", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	service.AssertExpectations(t)
}
