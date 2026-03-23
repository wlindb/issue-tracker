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
	trackerdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/project"
)

// mockProjectService implements the projectServicer interface for testing.
type mockProjectService struct {
	mock.Mock
}

func (m *mockProjectService) Create(ctx context.Context, id uuid.UUID, ownerID uuid.UUID, name string, description *string) (*trackerdomain.Project, error) {
	args := m.Called(ctx, id, ownerID, name, description)
	if p, ok := args.Get(0).(*trackerdomain.Project); ok {
		return p, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockProjectService) List(ctx context.Context, query trackerdomain.ListProjectQuery) (trackerdomain.Projects, error) {
	args := m.Called(ctx, query)
	projects, _ := args.Get(0).(trackerdomain.Projects)
	return projects, args.Error(1)
}

// newTestServer builds a minimal Echo server wired to the given project service.
func newTestServer(t *testing.T, service api.ProjectService) *echo.Echo {
	t.Helper()
	e := echo.New()
	h := &api.Handler{
		ProjectHandler: api.NewProjectHandler(service),
	}
	strict := model.NewStrictHandler(h, nil)
	model.RegisterHandlersWithBaseURL(e, strict, "/api/v1")
	return e
}

func Test_CreateProject_ValidBody_Returns201(t *testing.T) {
	service := &mockProjectService{}
	now := time.Now().UTC()
	ownerID := uuid.New()

	service.On("Create", mock.Anything, mock.Anything, ownerID, "Acme", (*string)(nil)).
		Return(&trackerdomain.Project{
			ID:        uuid.New(),
			Name:      "Acme",
			OwnerID:   ownerID,
			CreatedAt: now,
			UpdatedAt: now,
		}, nil)

	e := newTestServer(t, service)
	e.Use(injectUser(ownerID))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/projects", strings.NewReader(`{"name":"Acme"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var got model.Project
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	assert.Equal(t, "Acme", got.Name)
	service.AssertExpectations(t)
}

func Test_CreateProject_EmptyName_BadRequest(t *testing.T) {
	service := &mockProjectService{}

	e := newTestServer(t, service)
	e.Use(injectUser(uuid.New()))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/projects", strings.NewReader(`{}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	service.AssertNotCalled(t, "Create")
}

func Test_CreateProject_ServiceError_InternalServerError(t *testing.T) {
	service := &mockProjectService{}
	ownerID := uuid.New()

	service.On("Create", mock.Anything, mock.Anything, ownerID, "Acme", (*string)(nil)).
		Return(nil, errors.New("db down"))

	e := newTestServer(t, service)
	e.Use(injectUser(ownerID))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/projects", strings.NewReader(`{"name":"Acme"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	service.AssertExpectations(t)
}

func Test_CreateProject_MissingUserID_InternalServerError(t *testing.T) {
	service := &mockProjectService{}

	e := newTestServer(t, service) // no injectUser — userID absent from context

	req := httptest.NewRequest(http.MethodPost, "/api/v1/projects", strings.NewReader(`{"name":"Acme"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	service.AssertNotCalled(t, "Create")
}

func Test_ListProjects_ProjectsExist_Returns200(t *testing.T) {
	service := &mockProjectService{}
	ownerID := uuid.New()
	now := time.Now().UTC()

	service.On("List", mock.Anything, trackerdomain.ListProjectQuery{}).
		Return(trackerdomain.Projects{
			Items: []trackerdomain.Project{
				{
					ID:        uuid.New(),
					Name:      "Alpha",
					OwnerID:   ownerID,
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
		}, nil)

	e := newTestServer(t, service)
	e.Use(injectUser(ownerID))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/projects", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got model.ProjectPage
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Len(t, got.Items, 1)
	assert.Equal(t, "Alpha", got.Items[0].Name)
	assert.Equal(t, ownerID, got.Items[0].OwnerId)
	service.AssertExpectations(t)
}

func Test_ListProjects_EmptyList_Returns200(t *testing.T) {
	service := &mockProjectService{}
	ownerID := uuid.New()

	service.On("List", mock.Anything, trackerdomain.ListProjectQuery{}).
		Return(trackerdomain.Projects{Items: []trackerdomain.Project{}}, nil)

	e := newTestServer(t, service)
	e.Use(injectUser(ownerID))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/projects", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got model.ProjectPage
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	assert.Empty(t, got.Items)
	service.AssertExpectations(t)
}

func Test_ListProjects_NoUserID_Returns401(t *testing.T) {
	service := &mockProjectService{}

	e := newTestServer(t, service) // no injectUser — userID absent from context

	req := httptest.NewRequest(http.MethodGet, "/api/v1/projects", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
	service.AssertNotCalled(t, "List")
}

func Test_ListProjects_ServiceError_Returns500(t *testing.T) {
	service := &mockProjectService{}
	ownerID := uuid.New()

	service.On("List", mock.Anything, trackerdomain.ListProjectQuery{}).
		Return(trackerdomain.Projects{}, errors.New("db down"))

	e := newTestServer(t, service)
	e.Use(injectUser(ownerID))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/projects", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	service.AssertExpectations(t)
}

func Test_ListProjects_InvalidLimitParam_Returns400(t *testing.T) {
	service := &mockProjectService{}
	ownerID := uuid.New()

	e := newTestServer(t, service)
	e.Use(injectUser(ownerID))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/projects?limit=notanumber", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	service.AssertNotCalled(t, "List")
}
