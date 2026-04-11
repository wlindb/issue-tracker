//go:build integration

package api_test

import (
	"encoding/json"
	"fmt"
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
	issuedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
	projectdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/project"
	workspacedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/workspace"
	tracker "github.com/wlindb/issue-tracker/internal/infrastructure/tracker"
)

// injectWorkspace returns an Echo middleware that injects both workspace and user IDs
// into the request context, simulating workspace membership middleware.
func injectWorkspace(workspaceID, userID uuid.UUID) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := api.WithUserID(c.Request().Context(), userID)
			ctx = api.WithWorkspaceID(ctx, workspaceID)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

type issueIntegrationFixture struct {
	workspaceID uuid.UUID
	userID      uuid.UUID
	projectID   uuid.UUID
}

func setupIssueFixture(t *testing.T) issueIntegrationFixture {
	t.Helper()

	wsRepo := tracker.NewWorkspaceRepository(testPool)
	workspaceID := uuid.New()
	ownerID := uuid.New()
	_, err := wsRepo.Create(t.Context(), workspacedomain.Workspace{ID: workspaceID, OwnerID: ownerID, Name: "IssueTest"})
	require.NoError(t, err)

	ctx := api.WithWorkspaceID(api.WithUserID(t.Context(), ownerID), workspaceID)
	projRepo := tracker.NewProjectRepository(testPool)
	projectID := uuid.New()
	_, err = projRepo.Create(ctx, projectID, ownerID, "IssueProject-"+projectID.String()[:8], nil)
	require.NoError(t, err)

	return issueIntegrationFixture{
		workspaceID: workspaceID,
		userID:      ownerID,
		projectID:   projectID,
	}
}

func newIssueIntegrationServer(t *testing.T, f issueIntegrationFixture) *echo.Echo {
	t.Helper()
	issueRepo := tracker.NewIssueRepository(testPool)
	issueSvc := issuedomain.NewIssueService(issueRepo)
	projRepo := tracker.NewProjectRepository(testPool)
	projSvc := projectdomain.NewProjectService(projRepo)

	h := &api.Handler{
		IssueHandler:   api.NewIssueHandler(issueSvc),
		ProjectHandler: api.NewProjectHandler(projSvc),
	}
	e := echo.New()
	e.Use(injectWorkspace(f.workspaceID, f.userID))
	strict := model.NewStrictHandler(h, nil)
	model.RegisterHandlersWithBaseURL(e, strict, "/api/v1")
	return e
}

func createIssueViaAPI(t *testing.T, e *echo.Echo, f issueIntegrationFixture) model.Issue {
	t.Helper()
	body := fmt.Sprintf(`{"projectId":"%s","title":"Test Issue","status":"backlog","priority":"none"}`, f.projectID)
	req := httptest.NewRequest(http.MethodPost, wsPath("/issues"), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)
	var created model.Issue
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &created))
	return created
}

// — UpdateIssueDescription integration —

func Test_UpdateIssueDescription_Integration_Returns200(t *testing.T) {
	f := setupIssueFixture(t)
	e := newIssueIntegrationServer(t, f)
	created := createIssueViaAPI(t, e, f)

	body := `{"description":"updated via integration test"}`
	req := httptest.NewRequest(http.MethodPut, wsPath("/issues/"+created.Id.String()+"/description"), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var actual model.Issue
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actual))
	require.NotNil(t, actual.Description)
	assert.Equal(t, "updated via integration test", *actual.Description)
}

// — UpdateIssuePriority integration —

func Test_UpdateIssuePriority_Integration_Returns200(t *testing.T) {
	f := setupIssueFixture(t)
	e := newIssueIntegrationServer(t, f)
	created := createIssueViaAPI(t, e, f)

	body := `{"priority":"high"}`
	req := httptest.NewRequest(http.MethodPut, wsPath("/issues/"+created.Id.String()+"/priority"), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var actual model.Issue
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actual))
	assert.Equal(t, model.High, actual.Priority)
}

// — UpdateIssueStatus integration —

func Test_UpdateIssueStatus_Integration_Returns200(t *testing.T) {
	f := setupIssueFixture(t)
	e := newIssueIntegrationServer(t, f)
	created := createIssueViaAPI(t, e, f)

	body := `{"status":"in_progress"}`
	req := httptest.NewRequest(http.MethodPut, wsPath("/issues/"+created.Id.String()+"/status"), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var actual model.Issue
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actual))
	assert.Equal(t, model.InProgress, actual.Status)
}

// — UpdateIssueAssignee integration —

func Test_UpdateIssueAssignee_Integration_Returns200(t *testing.T) {
	f := setupIssueFixture(t)
	e := newIssueIntegrationServer(t, f)
	created := createIssueViaAPI(t, e, f)

	assigneeID := uuid.New()
	body := fmt.Sprintf(`{"assigneeId":"%s"}`, assigneeID)
	req := httptest.NewRequest(http.MethodPut, wsPath("/issues/"+created.Id.String()+"/assigneeId"), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var actual model.Issue
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actual))
	require.NotNil(t, actual.AssigneeId)
	assert.Equal(t, assigneeID, *actual.AssigneeId)
}
