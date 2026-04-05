//go:build !integration

package api_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/wlindb/issue-tracker/internal/api"
	"github.com/wlindb/issue-tracker/internal/api/model"
	workspacedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/workspace"
)

type mockWorkspaceService struct {
	mock.Mock
}

func (m *mockWorkspaceService) Create(ctx context.Context, ownerID uuid.UUID, name string) (*workspacedomain.Workspace, error) {
	args := m.Called(ctx, ownerID, name)
	if w, ok := args.Get(0).(*workspacedomain.Workspace); ok {
		return w, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockWorkspaceService) Get(ctx context.Context, id uuid.UUID) (*workspacedomain.Workspace, error) {
	args := m.Called(ctx, id)
	if w, ok := args.Get(0).(*workspacedomain.Workspace); ok {
		return w, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockWorkspaceService) List(ctx context.Context, userID uuid.UUID) ([]workspacedomain.Workspace, error) {
	args := m.Called(ctx, userID)
	workspaces, _ := args.Get(0).([]workspacedomain.Workspace)
	return workspaces, args.Error(1)
}

func newWorkspaceTestServer(t *testing.T, service api.WorkspaceService) *echo.Echo {
	t.Helper()
	e := echo.New()
	h := &api.Handler{
		WorkspaceHandler: api.NewWorkspaceHandler(service),
	}
	strict := model.NewStrictHandler(h, nil)
	model.RegisterHandlersWithBaseURL(e, strict, "/api/v1")
	return e
}

func Test_CreateWorkspace_ValidBody_Returns201(t *testing.T) {
	service := &mockWorkspaceService{}
	ownerID := uuid.New()
	now := time.Now().UTC()

	service.On("Create", mock.Anything, ownerID, "Acme").
		Return(&workspacedomain.Workspace{
			ID:        uuid.New(),
			Name:      "Acme",
			OwnerID:   ownerID,
			CreatedAt: now,
			UpdatedAt: now,
		}, nil)

	e := newWorkspaceTestServer(t, service)
	e.Use(injectUser(ownerID))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/workspaces", strings.NewReader(`{"name":"Acme"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var got model.Workspace
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	assert.Equal(t, "Acme", got.Name)
	service.AssertExpectations(t)
}

func Test_CreateWorkspace_EmptyName_Returns400(t *testing.T) {
	service := &mockWorkspaceService{}

	e := newWorkspaceTestServer(t, service)
	e.Use(injectUser(uuid.New()))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/workspaces", strings.NewReader(`{}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	service.AssertNotCalled(t, "Create")
}

func Test_CreateWorkspace_MissingUserID_ReturnsInternalServerError(t *testing.T) {
	service := &mockWorkspaceService{}

	e := newWorkspaceTestServer(t, service) // no injectUser

	req := httptest.NewRequest(http.MethodPost, "/api/v1/workspaces", strings.NewReader(`{"name":"Acme"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	service.AssertNotCalled(t, "Create")
}

func Test_CreateWorkspace_ServiceError_ReturnsInternalServerError(t *testing.T) {
	service := &mockWorkspaceService{}
	ownerID := uuid.New()

	service.On("Create", mock.Anything, ownerID, "Acme").
		Return(nil, errors.New("db down"))

	e := newWorkspaceTestServer(t, service)
	e.Use(injectUser(ownerID))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/workspaces", strings.NewReader(`{"name":"Acme"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	service.AssertExpectations(t)
}

func Test_GetWorkspace_ExistingWorkspace_Returns200(t *testing.T) {
	service := &mockWorkspaceService{}
	workspaceID := uuid.New()
	now := time.Now().UTC()

	service.On("Get", mock.Anything, workspaceID).
		Return(&workspacedomain.Workspace{
			ID:        workspaceID,
			Name:      "Acme",
			OwnerID:   uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
		}, nil)

	e := newWorkspaceTestServer(t, service)
	e.Use(injectUser(uuid.New()))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/workspaces/"+workspaceID.String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got model.Workspace
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	assert.Equal(t, "Acme", got.Name)
	service.AssertExpectations(t)
}

func Test_GetWorkspace_NotFound_Returns404(t *testing.T) {
	service := &mockWorkspaceService{}
	workspaceID := uuid.New()

	service.On("Get", mock.Anything, workspaceID).
		Return(nil, workspacedomain.ErrWorkspaceNotFound)

	e := newWorkspaceTestServer(t, service)
	e.Use(injectUser(uuid.New()))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/workspaces/"+workspaceID.String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
	service.AssertExpectations(t)
}

func Test_ListWorkspaces_WorkspacesExist_Returns200(t *testing.T) {
	service := &mockWorkspaceService{}
	userID := uuid.New()
	now := time.Now().UTC()

	service.On("List", mock.Anything, userID).
		Return([]workspacedomain.Workspace{
			{ID: uuid.New(), Name: "Acme", OwnerID: userID, CreatedAt: now, UpdatedAt: now},
		}, nil)

	e := newWorkspaceTestServer(t, service)
	e.Use(injectUser(userID))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/workspaces", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got model.WorkspacePage
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Len(t, got.Items, 1)
	assert.Equal(t, "Acme", got.Items[0].Name)
	service.AssertExpectations(t)
}

func Test_ListWorkspaces_NoUserID_Returns401(t *testing.T) {
	service := &mockWorkspaceService{}

	e := newWorkspaceTestServer(t, service) // no injectUser

	req := httptest.NewRequest(http.MethodGet, "/api/v1/workspaces", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
	service.AssertNotCalled(t, "List")
}

func Test_ListWorkspaces_ServiceError_Returns500(t *testing.T) {
	service := &mockWorkspaceService{}
	userID := uuid.New()

	service.On("List", mock.Anything, userID).Return(nil, errors.New("db down"))

	e := newWorkspaceTestServer(t, service)
	e.Use(injectUser(userID))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/workspaces", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	service.AssertExpectations(t)
}
