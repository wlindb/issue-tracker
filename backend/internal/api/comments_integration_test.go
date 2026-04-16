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
	commentdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/comment"
	issuedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
	projectdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/project"
	workspacedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/workspace"
	tracker "github.com/wlindb/issue-tracker/internal/infrastructure/tracker"
)

type commentIntegrationFixture struct {
	workspaceID uuid.UUID
	userID      uuid.UUID
	projectID   uuid.UUID
	issueID     uuid.UUID
}

func setupCommentFixture(t *testing.T) commentIntegrationFixture {
	t.Helper()

	wsRepo := tracker.NewWorkspaceRepository(testPool)
	workspaceID := uuid.New()
	ownerID := uuid.New()
	_, err := wsRepo.Create(t.Context(), workspacedomain.Workspace{ID: workspaceID, OwnerID: ownerID, Name: "CommentTest"})
	require.NoError(t, err)

	ctx := api.WithWorkspaceID(api.WithUserID(t.Context(), ownerID), workspaceID)

	projRepo := tracker.NewProjectRepository(testPool)
	projectID := uuid.New()
	project, err := projectdomain.New(projectID, strings.ToLower(projectID.String()[:8]), "CommentProject", nil, ownerID)
	require.NoError(t, err)
	_, err = projRepo.Create(ctx, project)
	require.NoError(t, err)

	issueRepo := tracker.NewIssueRepository(testPool)
	issueID := uuid.New()
	issue := issuedomain.Issue{
		ID:         issueID,
		Identifier: "comment-test-" + issueID.String()[:8],
		Title:      "Comment Test Issue",
		Status:     issuedomain.StatusBacklog,
		Priority:   issuedomain.PriorityNone,
		Labels:     []string{},
		ProjectID:  projectID,
		ReporterID: ownerID,
	}
	_, err = issueRepo.CreateIssue(ctx, issue)
	require.NoError(t, err)

	return commentIntegrationFixture{
		workspaceID: workspaceID,
		userID:      ownerID,
		projectID:   projectID,
		issueID:     issueID,
	}
}

func newCommentIntegrationServer(t *testing.T, f commentIntegrationFixture) *echo.Echo {
	t.Helper()
	commentRepo := tracker.NewCommentRepository(testPool)
	commentSvc := commentdomain.NewCommentService(commentRepo)

	h := &api.Handler{
		CommentHandler: api.NewCommentHandler(commentSvc),
	}
	e := echo.New()
	e.Use(injectWorkspace(f.workspaceID, f.userID))
	strict := model.NewStrictHandler(h, nil)
	model.RegisterHandlersWithBaseURL(e, strict, "/api/v1")
	return e
}

func Test_CreateComment_ValidBody_Returns201(t *testing.T) {
	f := setupCommentFixture(t)
	e := newCommentIntegrationServer(t, f)

	body := `{"body":"Hello from integration test"}`
	path := fmt.Sprintf("/api/v1/workspaces/%s/issues/%s/comments", f.workspaceID, f.issueID)
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var actual model.Comment
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actual))
	assert.Equal(t, "Hello from integration test", actual.Body)
	assert.Equal(t, f.issueID, actual.IssueId)
	assert.Equal(t, f.userID, actual.AuthorId)
	assert.NotEqual(t, uuid.Nil, actual.Id)
}

func Test_ListComments_CommentsExist_Returns200(t *testing.T) {
	f := setupCommentFixture(t)
	e := newCommentIntegrationServer(t, f)

	// Create a comment first
	createBody := `{"body":"Comment for listing"}`
	createPath := fmt.Sprintf("/api/v1/workspaces/%s/issues/%s/comments", f.workspaceID, f.issueID)
	createReq := httptest.NewRequest(http.MethodPost, createPath, strings.NewReader(createBody))
	createReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	createRec := httptest.NewRecorder()
	e.ServeHTTP(createRec, createReq)
	require.Equal(t, http.StatusCreated, createRec.Code)

	// List comments
	listPath := fmt.Sprintf("/api/v1/workspaces/%s/issues/%s/comments", f.workspaceID, f.issueID)
	listReq := httptest.NewRequest(http.MethodGet, listPath, nil)
	listRec := httptest.NewRecorder()
	e.ServeHTTP(listRec, listReq)

	require.Equal(t, http.StatusOK, listRec.Code)
	var actual model.CommentPage
	require.NoError(t, json.Unmarshal(listRec.Body.Bytes(), &actual))
	require.NotEmpty(t, actual.Items)
	assert.Equal(t, "Comment for listing", actual.Items[0].Body)
}
