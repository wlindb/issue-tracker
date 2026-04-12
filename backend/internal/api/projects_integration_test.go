//go:build integration

package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wlindb/issue-tracker/internal/api"
	"github.com/wlindb/issue-tracker/internal/api/model"
	projectdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/project"
	tracker "github.com/wlindb/issue-tracker/internal/infrastructure/tracker"
)

func newProjectIntegrationServer(t *testing.T) *echo.Echo {
	t.Helper()
	repository := tracker.NewProjectRepository(testPool)
	service := projectdomain.NewProjectService(repository)
	handler := api.NewProjectHandler(service)
	h := &api.Handler{ProjectHandler: handler}
	e := echo.New()
	strict := model.NewStrictHandler(h, nil)
	model.RegisterHandlersWithBaseURL(e, strict, "/api/v1")
	return e
}

func Test_GetProject_ValidRequest_Returns200(t *testing.T) {
	ownerID := uuid.New()

	// Create a workspace first so the FK constraint on workspace_id is satisfied.
	workspaceServer := newWorkspaceIntegrationServer(t)
	workspaceServer.Use(injectUser(ownerID))
	createWorkspaceReq := httptest.NewRequest(http.MethodPost, "/api/v1/workspaces", strings.NewReader(`{"name":"Test Workspace"}`))
	createWorkspaceReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	createWorkspaceRec := httptest.NewRecorder()
	workspaceServer.ServeHTTP(createWorkspaceRec, createWorkspaceReq)
	require.Equal(t, http.StatusCreated, createWorkspaceRec.Code)
	var createdWorkspace model.Workspace
	require.NoError(t, json.Unmarshal(createWorkspaceRec.Body.Bytes(), &createdWorkspace))
	workspaceID := createdWorkspace.Id

	e := newProjectIntegrationServer(t)
	e.Use(injectUser(ownerID))
	e.Use(injectWorkspace(workspaceID, ownerID))

	// Create the project first.
	createReq := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/workspaces/"+workspaceID.String()+"/projects",
		strings.NewReader(`{"name":"Get Test Project"}`),
	)
	createReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	createRec := httptest.NewRecorder()
	e.ServeHTTP(createRec, createReq)
	require.Equal(t, http.StatusCreated, createRec.Code)
	var created model.Project
	require.NoError(t, json.Unmarshal(createRec.Body.Bytes(), &created))

	// Now fetch the project by ID.
	getReq := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/workspaces/"+workspaceID.String()+"/projects/"+created.Id.String(),
		nil,
	)
	getRec := httptest.NewRecorder()
	e.ServeHTTP(getRec, getReq)

	require.Equal(t, http.StatusOK, getRec.Code)
	var actual model.Project
	require.NoError(t, json.Unmarshal(getRec.Body.Bytes(), &actual))
	assert.Equal(t, created.Id, actual.Id)
	assert.Equal(t, "Get Test Project", actual.Name)
	assert.Equal(t, ownerID, actual.OwnerId)
}
	ownerID := uuid.New()

	// Create a workspace first so the FK constraint on workspace_id is satisfied.
	workspaceServer := newWorkspaceIntegrationServer(t)
	workspaceServer.Use(injectUser(ownerID))
	createWorkspaceReq := httptest.NewRequest(http.MethodPost, "/api/v1/workspaces", strings.NewReader(`{"name":"Test Workspace"}`))
	createWorkspaceReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	createWorkspaceRec := httptest.NewRecorder()
	workspaceServer.ServeHTTP(createWorkspaceRec, createWorkspaceReq)
	require.Equal(t, http.StatusCreated, createWorkspaceRec.Code)
	var createdWorkspace model.Workspace
	require.NoError(t, json.Unmarshal(createWorkspaceRec.Body.Bytes(), &createdWorkspace))
	workspaceID := createdWorkspace.Id

	e := newProjectIntegrationServer(t)
	e.Use(injectUser(ownerID))
	e.Use(injectWorkspace(workspaceID, ownerID))

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/workspaces/"+workspaceID.String()+"/projects",
		strings.NewReader(`{"name":"My Project"}`),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var actual model.Project
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actual))
	assert.Equal(t, "My Project", actual.Name)
	assert.NotEqual(t, uuid.Nil, actual.Id)
	assert.Equal(t, ownerID, actual.OwnerId)
	assert.NotEmpty(t, actual.Identifier)
}
