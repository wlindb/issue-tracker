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
	issuedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
)

// mockIssueService implements api.IssueService for testing.
type mockIssueService struct {
	mock.Mock
}

func (m *mockIssueService) ListIssues(ctx context.Context, projectID uuid.UUID, query issuedomain.ListIssueQuery) (issuedomain.IssuePage, error) {
	args := m.Called(ctx, projectID, query)
	if page, ok := args.Get(0).(issuedomain.IssuePage); ok {
		return page, args.Error(1)
	}
	return issuedomain.IssuePage{}, args.Error(1)
}

func (m *mockIssueService) CreateIssue(ctx context.Context, command issuedomain.CreateIssueCommand) (*issuedomain.Issue, error) {
	args := m.Called(ctx, command)
	if issue, ok := args.Get(0).(*issuedomain.Issue); ok {
		return issue, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockIssueService) GetIssue(ctx context.Context, issueID uuid.UUID) (*issuedomain.Issue, error) {
	args := m.Called(ctx, issueID)
	if issue, ok := args.Get(0).(*issuedomain.Issue); ok {
		return issue, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockIssueService) UpdateIssuePriority(ctx context.Context, issueID uuid.UUID, priority issuedomain.Priority) (*issuedomain.Issue, error) {
	args := m.Called(ctx, issueID, priority)
	if issue, ok := args.Get(0).(*issuedomain.Issue); ok {
		return issue, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockIssueService) UpdateIssueAssignee(ctx context.Context, issueID uuid.UUID, assigneeID *uuid.UUID) (*issuedomain.Issue, error) {
	args := m.Called(ctx, issueID, assigneeID)
	if issue, ok := args.Get(0).(*issuedomain.Issue); ok {
		return issue, args.Error(1)
	}
	return nil, args.Error(1)
}

// newIssueTestServer builds a minimal Echo server wired to the given issue service.
func newIssueTestServer(t *testing.T, service api.IssueService) *echo.Echo {
	t.Helper()
	e := echo.New()
	h := &api.Handler{
		IssueHandler: api.NewIssueHandler(service),
	}
	strict := model.NewStrictHandler(h, nil)
	model.RegisterHandlersWithBaseURL(e, strict, "/api/v1")
	return e
}

// — ListIssues —

func Test_ListIssues_IssuesExist_Returns200(t *testing.T) {
	service := &mockIssueService{}
	projectID := uuid.New()
	userID := uuid.New()
	now := time.Now().UTC()

	service.On("ListIssues", mock.Anything, projectID, mock.Anything).
		Return(issuedomain.IssuePage{
			Items: []issuedomain.Issue{
				{
					ID:         uuid.New(),
					Identifier: "PROJ-1",
					Title:      "Fix login bug",
					Status:     issuedomain.StatusBacklog,
					Priority:   issuedomain.PriorityNone,
					Labels:     []string{},
					ProjectID:  projectID,
					ReporterID: userID,
					CreatedAt:  now,
					UpdatedAt:  now,
				},
			},
			NextCursor: nil,
		}, nil)

	e := newIssueTestServer(t, service)
	e.Use(injectUser(userID))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/issues?project_id="+projectID.String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var actual model.IssuePage
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actual))
	require.Len(t, actual.Items, 1)
	assert.Equal(t, "Fix login bug", actual.Items[0].Title)
	service.AssertExpectations(t)
}

func Test_ListIssues_EmptyList_Returns200(t *testing.T) {
	service := &mockIssueService{}
	projectID := uuid.New()

	service.On("ListIssues", mock.Anything, projectID, mock.Anything).
		Return(issuedomain.IssuePage{Items: []issuedomain.Issue{}, NextCursor: nil}, nil)

	e := newIssueTestServer(t, service)
	e.Use(injectUser(uuid.New()))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/issues?project_id="+projectID.String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var actual model.IssuePage
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actual))
	assert.Empty(t, actual.Items)
	service.AssertExpectations(t)
}

func Test_ListIssues_NoUserID_Returns401(t *testing.T) {
	service := &mockIssueService{}
	projectID := uuid.New()

	e := newIssueTestServer(t, service) // no injectUser — userID absent from context

	req := httptest.NewRequest(http.MethodGet, "/api/v1/issues?project_id="+projectID.String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
	service.AssertNotCalled(t, "ListIssues")
}

func Test_ListIssues_ProjectNotFound_Returns404(t *testing.T) {
	service := &mockIssueService{}
	projectID := uuid.New()

	service.On("ListIssues", mock.Anything, projectID, mock.Anything).
		Return(issuedomain.IssuePage{}, api.ErrIssueProjectNotFound)

	e := newIssueTestServer(t, service)
	e.Use(injectUser(uuid.New()))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/issues?project_id="+projectID.String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
	service.AssertExpectations(t)
}

func Test_ListIssues_ServiceError_Returns500(t *testing.T) {
	service := &mockIssueService{}
	projectID := uuid.New()

	service.On("ListIssues", mock.Anything, projectID, mock.Anything).
		Return(issuedomain.IssuePage{}, errors.New("db down"))

	e := newIssueTestServer(t, service)
	e.Use(injectUser(uuid.New()))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/issues?project_id="+projectID.String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	service.AssertExpectations(t)
}

func Test_ListIssues_InvalidLimitParam_Returns400(t *testing.T) {
	service := &mockIssueService{}
	projectID := uuid.New()

	e := newIssueTestServer(t, service)
	e.Use(injectUser(uuid.New()))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/issues?project_id="+projectID.String()+"&limit=notanumber", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	service.AssertNotCalled(t, "ListIssues")
}

// — CreateIssue —

func Test_CreateIssue_ValidBody_Returns201(t *testing.T) {
	service := &mockIssueService{}
	projectID := uuid.New()
	userID := uuid.New()
	now := time.Now().UTC()

	service.On("CreateIssue", mock.Anything, mock.Anything).
		Return(&issuedomain.Issue{
			ID:         uuid.New(),
			Identifier: "PROJ-1",
			Title:      "New feature",
			Status:     issuedomain.StatusTodo,
			Priority:   issuedomain.PriorityMedium,
			Labels:     []string{},
			ProjectID:  projectID,
			ReporterID: userID,
			CreatedAt:  now,
			UpdatedAt:  now,
		}, nil)

	e := newIssueTestServer(t, service)
	e.Use(injectUser(userID))

	body := `{"projectId":"` + projectID.String() + `","title":"New feature","status":"todo","priority":"medium"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/issues", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var actual model.Issue
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actual))
	assert.Equal(t, "New feature", actual.Title)
	service.AssertExpectations(t)
}

func Test_CreateIssue_MissingProjectID_Returns400(t *testing.T) {
	service := &mockIssueService{}

	e := newIssueTestServer(t, service)
	e.Use(injectUser(uuid.New()))

	body := `{"title":"New feature","status":"todo","priority":"medium"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/issues", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	service.AssertNotCalled(t, "CreateIssue")
}

func Test_CreateIssue_MissingTitle_Returns400(t *testing.T) {
	service := &mockIssueService{}
	projectID := uuid.New()

	e := newIssueTestServer(t, service)
	e.Use(injectUser(uuid.New()))

	body := `{"projectId":"` + projectID.String() + `","status":"backlog","priority":"none"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/issues", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	service.AssertNotCalled(t, "CreateIssue")
}

func Test_CreateIssue_NoUserID_Returns401(t *testing.T) {
	service := &mockIssueService{}
	projectID := uuid.New()

	e := newIssueTestServer(t, service) // no injectUser — userID absent from context

	body := `{"projectId":"` + projectID.String() + `","title":"New feature","status":"todo","priority":"medium"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/issues", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
	service.AssertNotCalled(t, "CreateIssue")
}

func Test_CreateIssue_InvalidProjectID_Returns422(t *testing.T) {
	service := &mockIssueService{}
	projectID := uuid.New()
	userID := uuid.New()

	service.On("CreateIssue", mock.Anything, mock.Anything).
		Return(nil, api.ErrIssueUnprocessable)

	e := newIssueTestServer(t, service)
	e.Use(injectUser(userID))

	body := `{"projectId":"` + projectID.String() + `","title":"New feature","status":"todo","priority":"medium"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/issues", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	service.AssertExpectations(t)
}

func Test_CreateIssue_UnprocessableInput_Returns422(t *testing.T) {
	service := &mockIssueService{}
	projectID := uuid.New()
	userID := uuid.New()

	service.On("CreateIssue", mock.Anything, mock.Anything).
		Return(nil, api.ErrIssueUnprocessable)

	e := newIssueTestServer(t, service)
	e.Use(injectUser(userID))

	body := `{"projectId":"` + projectID.String() + `","title":"New feature","status":"todo","priority":"medium"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/issues", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	service.AssertExpectations(t)
}

func Test_CreateIssue_ServiceError_Returns500(t *testing.T) {
	service := &mockIssueService{}
	projectID := uuid.New()
	userID := uuid.New()

	service.On("CreateIssue", mock.Anything, mock.Anything).
		Return(nil, errors.New("db down"))

	e := newIssueTestServer(t, service)
	e.Use(injectUser(userID))

	body := `{"projectId":"` + projectID.String() + `","title":"New feature","status":"todo","priority":"medium"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/issues", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	service.AssertExpectations(t)
}

// — SearchIssues —

func Test_SearchIssues_Returns501(t *testing.T) {
	service := &mockIssueService{}

	e := newIssueTestServer(t, service)
	e.Use(injectUser(uuid.New()))

	body := `{"query":"login bug"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/issues/search", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotImplemented, rec.Code)
	var actual model.Error
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actual))
	assert.Equal(t, "not_implemented", actual.Code)
	service.AssertNotCalled(t, "ListIssues")
	service.AssertNotCalled(t, "CreateIssue")
}

func Test_SearchIssues_MissingBody_Returns400(t *testing.T) {
	service := &mockIssueService{}

	e := newIssueTestServer(t, service)
	e.Use(injectUser(uuid.New()))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/issues/search", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	service.AssertNotCalled(t, "ListIssues")
	service.AssertNotCalled(t, "CreateIssue")
}

func Test_SearchIssues_MissingQuery_Returns400(t *testing.T) {
	service := &mockIssueService{}

	e := newIssueTestServer(t, service)
	e.Use(injectUser(uuid.New()))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/issues/search", strings.NewReader(`{}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	service.AssertNotCalled(t, "ListIssues")
	service.AssertNotCalled(t, "CreateIssue")
}

func Test_SearchIssues_EmptyQuery_Returns400(t *testing.T) {
	service := &mockIssueService{}

	e := newIssueTestServer(t, service)
	e.Use(injectUser(uuid.New()))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/issues/search", strings.NewReader(`{"query":""}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	service.AssertNotCalled(t, "ListIssues")
	service.AssertNotCalled(t, "CreateIssue")
}

// — UpdateIssuePriority —

func Test_UpdateIssuePriority_ValidRequest_Returns200(t *testing.T) {
	service := &mockIssueService{}
	issueID := uuid.New()
	updatedIssue := &issuedomain.Issue{
		ID:         issueID,
		Identifier: "ISS-1",
		Title:      "Test issue",
		Status:     issuedomain.StatusTodo,
		Priority:   issuedomain.PriorityHigh,
		Labels:     []string{},
		ProjectID:  uuid.New(),
		ReporterID: uuid.New(),
	}

	service.On("UpdateIssuePriority", mock.Anything, issueID, issuedomain.PriorityHigh).
		Return(updatedIssue, nil)

	e := newIssueTestServer(t, service)
	e.Use(injectUser(uuid.New()))

	body := `{"priority":"high"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/issues/"+issueID.String()+"/priority", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	service.AssertExpectations(t)
}

func Test_UpdateIssuePriority_NoUserID_Returns401(t *testing.T) {
	service := &mockIssueService{}
	issueID := uuid.New()

	e := newIssueTestServer(t, service) // no injectUser — userID absent from context

	body := `{"priority":"high"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/issues/"+issueID.String()+"/priority", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
	service.AssertNotCalled(t, "UpdateIssuePriority")
}

func Test_UpdateIssuePriority_InvalidPriority_Returns400(t *testing.T) {
	service := &mockIssueService{}
	issueID := uuid.New()

	e := newIssueTestServer(t, service)
	e.Use(injectUser(uuid.New()))

	body := `{"priority":"not-a-priority"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/issues/"+issueID.String()+"/priority", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	service.AssertNotCalled(t, "UpdateIssuePriority")
}

func Test_UpdateIssuePriority_IssueNotFound_Returns422(t *testing.T) {
	service := &mockIssueService{}
	issueID := uuid.New()

	service.On("UpdateIssuePriority", mock.Anything, issueID, issuedomain.PriorityLow).
		Return((*issuedomain.Issue)(nil), api.ErrIssueNotFound)

	e := newIssueTestServer(t, service)
	e.Use(injectUser(uuid.New()))

	body := `{"priority":"low"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/issues/"+issueID.String()+"/priority", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	service.AssertExpectations(t)
}

func Test_UpdateIssuePriority_ServiceError_Returns500(t *testing.T) {
	service := &mockIssueService{}
	issueID := uuid.New()

	service.On("UpdateIssuePriority", mock.Anything, issueID, issuedomain.PriorityMedium).
		Return((*issuedomain.Issue)(nil), errors.New("db down"))

	e := newIssueTestServer(t, service)
	e.Use(injectUser(uuid.New()))

	body := `{"priority":"medium"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/issues/"+issueID.String()+"/priority", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	service.AssertExpectations(t)
}

// — GetIssue —

func Test_GetIssue_NoUserID_Returns401(t *testing.T) {
	service := &mockIssueService{}
	issueID := uuid.New()

	e := newIssueTestServer(t, service) // no injectUser — userID absent from context

	req := httptest.NewRequest(http.MethodGet, "/api/v1/issues/"+issueID.String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
	service.AssertNotCalled(t, "GetIssue")
}

func Test_GetIssue_IssueNotFound_Returns404(t *testing.T) {
	service := &mockIssueService{}
	issueID := uuid.New()

	service.On("GetIssue", mock.Anything, issueID).
		Return(nil, api.ErrIssueNotFound)

	e := newIssueTestServer(t, service)
	e.Use(injectUser(uuid.New()))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/issues/"+issueID.String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
	var actual model.Error
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actual))
	assert.Equal(t, "not_found", actual.Code)
	service.AssertExpectations(t)
}

func Test_GetIssue_ServiceError_Returns500(t *testing.T) {
	service := &mockIssueService{}
	issueID := uuid.New()

	service.On("GetIssue", mock.Anything, issueID).
		Return(nil, errors.New("db down"))

	e := newIssueTestServer(t, service)
	e.Use(injectUser(uuid.New()))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/issues/"+issueID.String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	service.AssertExpectations(t)
}

func Test_GetIssue_IssueExists_Returns200(t *testing.T) {
	service := &mockIssueService{}
	issueID := uuid.New()
	projectID := uuid.New()
	reporterID := uuid.New()
	now := time.Now().UTC().Truncate(time.Second)

	service.On("GetIssue", mock.Anything, issueID).
		Return(&issuedomain.Issue{
			ID:         issueID,
			Identifier: "PROJ-1",
			Title:      "Fix login bug",
			Status:     issuedomain.StatusBacklog,
			Priority:   issuedomain.PriorityNone,
			Labels:     []string{},
			ProjectID:  projectID,
			ReporterID: reporterID,
			CreatedAt:  now,
			UpdatedAt:  now,
		}, nil)

	e := newIssueTestServer(t, service)
	e.Use(injectUser(uuid.New()))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/issues/"+issueID.String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var actual model.Issue
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actual))
	assert.Equal(t, issueID, actual.Id)
	assert.Equal(t, "PROJ-1", actual.Identifier)
	assert.Equal(t, "Fix login bug", actual.Title)
	assert.Equal(t, model.IssueStatus(issuedomain.StatusBacklog), actual.Status)
	assert.Equal(t, model.IssuePriority(issuedomain.PriorityNone), actual.Priority)
	assert.Equal(t, projectID, actual.ProjectId)
	assert.Equal(t, reporterID, actual.ReporterId)
	service.AssertExpectations(t)
}

// — UpdateIssueAssignee —

func Test_UpdateIssueAssignee_ValidAssignee_Returns200(t *testing.T) {
	service := &mockIssueService{}
	issueID := uuid.New()
	assigneeID := uuid.New()
	userID := uuid.New()
	now := time.Now().UTC()

	updatedIssue := &issuedomain.Issue{
		ID:         issueID,
		Identifier: "PROJ-1",
		Title:      "Fix login bug",
		Status:     issuedomain.StatusBacklog,
		Priority:   issuedomain.PriorityNone,
		Labels:     []string{},
		ProjectID:  uuid.New(),
		ReporterID: userID,
		AssigneeID: &assigneeID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	service.On("UpdateIssueAssignee", mock.Anything, issueID, &assigneeID).Return(updatedIssue, nil)

	e := newIssueTestServer(t, service)
	e.Use(injectUser(userID))

	body := `{"assigneeId":"` + assigneeID.String() + `"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/issues/"+issueID.String()+"/assigneeId", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var actual model.Issue
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actual))
	assert.Equal(t, "Fix login bug", actual.Title)
	service.AssertExpectations(t)
}

func Test_UpdateIssueAssignee_NullAssignee_Returns200(t *testing.T) {
	service := &mockIssueService{}
	issueID := uuid.New()
	userID := uuid.New()
	now := time.Now().UTC()

	updatedIssue := &issuedomain.Issue{
		ID:         issueID,
		Identifier: "PROJ-1",
		Title:      "Fix login bug",
		Status:     issuedomain.StatusBacklog,
		Priority:   issuedomain.PriorityNone,
		Labels:     []string{},
		ProjectID:  uuid.New(),
		ReporterID: userID,
		AssigneeID: nil,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	service.On("UpdateIssueAssignee", mock.Anything, issueID, (*uuid.UUID)(nil)).Return(updatedIssue, nil)

	e := newIssueTestServer(t, service)
	e.Use(injectUser(userID))

	body := `{"assigneeId":null}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/issues/"+issueID.String()+"/assigneeId", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var actual model.Issue
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actual))
	assert.Nil(t, actual.AssigneeId)
	service.AssertExpectations(t)
}

func Test_UpdateIssueAssignee_NoAuth_Returns401(t *testing.T) {
	service := &mockIssueService{}
	issueID := uuid.New()
	assigneeID := uuid.New()

	e := newIssueTestServer(t, service) // no injectUser — userID absent from context

	body := `{"assigneeId":"` + assigneeID.String() + `"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/issues/"+issueID.String()+"/assigneeId", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
	service.AssertNotCalled(t, "ListIssues")
	service.AssertNotCalled(t, "CreateIssue")
}

func Test_UpdateIssueAssignee_IssueNotFound_Returns422(t *testing.T) {
	service := &mockIssueService{}
	issueID := uuid.New()
	assigneeID := uuid.New()

	service.On("UpdateIssueAssignee", mock.Anything, issueID, &assigneeID).Return(nil, api.ErrIssueNotFound)

	e := newIssueTestServer(t, service)
	e.Use(injectUser(uuid.New()))

	body := `{"assigneeId":"` + assigneeID.String() + `"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/issues/"+issueID.String()+"/assigneeId", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	var actual model.Error
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actual))
	assert.Equal(t, "unprocessable_entity", actual.Code)
	service.AssertExpectations(t)
}

func Test_UpdateIssueAssignee_InvalidAssigneeIDFormat_Returns400(t *testing.T) {
	service := &mockIssueService{}
	issueID := uuid.New()

	e := newIssueTestServer(t, service)
	e.Use(injectUser(uuid.New()))

	body := `{"assigneeId":"not-a-uuid"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/issues/"+issueID.String()+"/assigneeId", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	service.AssertNotCalled(t, "ListIssues")
	service.AssertNotCalled(t, "CreateIssue")
}

func Test_UpdateIssueAssignee_InvalidIssueIDFormat_Returns400(t *testing.T) {
	service := &mockIssueService{}

	e := newIssueTestServer(t, service)
	e.Use(injectUser(uuid.New()))

	assigneeID := uuid.New()
	body := `{"assigneeId":"` + assigneeID.String() + `"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/issues/not-a-uuid/assigneeId", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	service.AssertNotCalled(t, "ListIssues")
	service.AssertNotCalled(t, "CreateIssue")
}
