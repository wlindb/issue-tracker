//go:build !integration

package api_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	commentdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/comment"
)

// mockCommentService implements CommentService for testing.
type mockCommentService struct {
	mock.Mock
}

func (m *mockCommentService) List(ctx context.Context, issueID uuid.UUID, query commentdomain.ListCommentQuery) (commentdomain.Comments, error) {
	args := m.Called(ctx, issueID, query)
	comments, _ := args.Get(0).(commentdomain.Comments)
	return comments, args.Error(1)
}

func (m *mockCommentService) Create(ctx context.Context, comment commentdomain.Comment) (*commentdomain.Comment, error) {
	args := m.Called(ctx, comment)
	if c, ok := args.Get(0).(*commentdomain.Comment); ok {
		return c, args.Error(1)
	}
	return nil, args.Error(1)
}

// newCommentTestServer builds a minimal Echo server wired to the given comment service.
func newCommentTestServer(t *testing.T, service api.CommentService) *echo.Echo {
	t.Helper()
	e := echo.New()
	e.HTTPErrorHandler = api.HTTPErrorHandler
	h := &api.Handler{
		ProjectHandler: api.NewProjectHandler(new(mockProjectService)),
		CommentHandler: api.NewCommentHandler(service),
	}
	strict := model.NewStrictHandler(h, nil)
	model.RegisterHandlersWithBaseURL(e, strict, "/api/v1")
	return e
}

// --- GET /issues/{issueId}/comments ---

func Test_ListComments_CommentsExist_Returns200(t *testing.T) {
	service := &mockCommentService{}
	issueID := uuid.New()
	authorID := uuid.New()
	now := time.Now().UTC()

	service.On("List", mock.Anything, issueID, commentdomain.NewListCommentQuery(nil, nil)).
		Return(commentdomain.Comments{
			Items: []commentdomain.Comment{
				{
					ID:        uuid.New(),
					Body:      "Hello world",
					AuthorID:  authorID,
					IssueID:   issueID,
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
		}, nil)

	e := newCommentTestServer(t, service)
	e.Use(injectUser(authorID))

	req := httptest.NewRequest(http.MethodGet, wsPath(fmt.Sprintf("/issues/%s/comments", issueID)), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got model.CommentPage
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Len(t, got.Items, 1)
	assert.Equal(t, "Hello world", got.Items[0].Body)
	assert.Equal(t, issueID, got.Items[0].IssueId)
	service.AssertExpectations(t)
}

func Test_ListComments_EmptyList_Returns200(t *testing.T) {
	service := &mockCommentService{}
	issueID := uuid.New()
	authorID := uuid.New()

	service.On("List", mock.Anything, issueID, commentdomain.NewListCommentQuery(nil, nil)).
		Return(commentdomain.Comments{Items: []commentdomain.Comment{}}, nil)

	e := newCommentTestServer(t, service)
	e.Use(injectUser(authorID))

	req := httptest.NewRequest(http.MethodGet, wsPath(fmt.Sprintf("/issues/%s/comments", issueID)), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got model.CommentPage
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	assert.Empty(t, got.Items)
	service.AssertExpectations(t)
}

func Test_ListComments_NoUserID_Returns401(t *testing.T) {
	service := &mockCommentService{}
	issueID := uuid.New()

	e := newCommentTestServer(t, service) // no injectUser — userID absent from context

	req := httptest.NewRequest(http.MethodGet, wsPath(fmt.Sprintf("/issues/%s/comments", issueID)), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
	service.AssertNotCalled(t, "List")
}

func Test_ListComments_IssueNotFound_Returns404(t *testing.T) {
	service := &mockCommentService{}
	issueID := uuid.New()
	authorID := uuid.New()

	service.On("List", mock.Anything, issueID, commentdomain.NewListCommentQuery(nil, nil)).
		Return(commentdomain.Comments{}, commentdomain.ErrIssueNotFound)

	e := newCommentTestServer(t, service)
	e.Use(injectUser(authorID))

	req := httptest.NewRequest(http.MethodGet, wsPath(fmt.Sprintf("/issues/%s/comments", issueID)), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
	service.AssertExpectations(t)
}

func Test_ListComments_ServiceError_Returns500(t *testing.T) {
	service := &mockCommentService{}
	issueID := uuid.New()
	authorID := uuid.New()

	service.On("List", mock.Anything, issueID, commentdomain.NewListCommentQuery(nil, nil)).
		Return(commentdomain.Comments{}, errors.New("db down"))

	e := newCommentTestServer(t, service)
	e.Use(injectUser(authorID))

	req := httptest.NewRequest(http.MethodGet, wsPath(fmt.Sprintf("/issues/%s/comments", issueID)), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	service.AssertExpectations(t)
}

// --- POST /issues/{issueId}/comments ---

func Test_CreateComment_ValidBody_Returns201(t *testing.T) {
	service := &mockCommentService{}
	issueID := uuid.New()
	authorID := uuid.New()
	now := time.Now().UTC()

	service.On("Create", mock.Anything, mock.MatchedBy(func(c commentdomain.Comment) bool {
		return c.IssueID == issueID && c.AuthorID == authorID && c.Body == "Hello world"
	})).
		Return(&commentdomain.Comment{
			ID:        uuid.New(),
			Body:      "Hello world",
			AuthorID:  authorID,
			IssueID:   issueID,
			CreatedAt: now,
			UpdatedAt: now,
		}, nil)

	e := newCommentTestServer(t, service)
	e.Use(injectUser(authorID))

	req := httptest.NewRequest(http.MethodPost, wsPath(fmt.Sprintf("/issues/%s/comments", issueID)), strings.NewReader(`{"body":"Hello world"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var got model.Comment
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	assert.Equal(t, "Hello world", got.Body)
	assert.Equal(t, issueID, got.IssueId)
	assert.Equal(t, authorID, got.AuthorId)
	service.AssertExpectations(t)
}

func Test_CreateComment_EmptyBody_Returns400(t *testing.T) {
	service := &mockCommentService{}
	issueID := uuid.New()

	e := newCommentTestServer(t, service)
	e.Use(injectUser(uuid.New()))

	req := httptest.NewRequest(http.MethodPost, wsPath(fmt.Sprintf("/issues/%s/comments", issueID)), strings.NewReader(`{"body":""}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	service.AssertNotCalled(t, "Create")
}

func Test_CreateComment_NoUserID_Returns401(t *testing.T) {
	service := &mockCommentService{}
	issueID := uuid.New()

	e := newCommentTestServer(t, service) // no injectUser — userID absent from context

	req := httptest.NewRequest(http.MethodPost, wsPath(fmt.Sprintf("/issues/%s/comments", issueID)), strings.NewReader(`{"body":"Hello"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
	service.AssertNotCalled(t, "Create")
}

func Test_CreateComment_IssueNotFound_Returns404(t *testing.T) {
	service := &mockCommentService{}
	issueID := uuid.New()
	authorID := uuid.New()

	service.On("Create", mock.Anything, mock.MatchedBy(func(c commentdomain.Comment) bool {
		return c.IssueID == issueID && c.AuthorID == authorID && c.Body == "Hello"
	})).
		Return(nil, commentdomain.ErrIssueNotFound)

	e := newCommentTestServer(t, service)
	e.Use(injectUser(authorID))

	req := httptest.NewRequest(http.MethodPost, wsPath(fmt.Sprintf("/issues/%s/comments", issueID)), strings.NewReader(`{"body":"Hello"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
	service.AssertExpectations(t)
}

func Test_CreateComment_ServiceError_Returns500(t *testing.T) {
	service := &mockCommentService{}
	issueID := uuid.New()
	authorID := uuid.New()

	service.On("Create", mock.Anything, mock.MatchedBy(func(c commentdomain.Comment) bool {
		return c.IssueID == issueID && c.AuthorID == authorID && c.Body == "Hello"
	})).
		Return(nil, errors.New("db down"))

	e := newCommentTestServer(t, service)
	e.Use(injectUser(authorID))

	req := httptest.NewRequest(http.MethodPost, wsPath(fmt.Sprintf("/issues/%s/comments", issueID)), strings.NewReader(`{"body":"Hello"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	service.AssertExpectations(t)
}
