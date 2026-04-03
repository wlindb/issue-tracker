//go:build !integration

package tracker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	issuedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
	trackerdb "github.com/wlindb/issue-tracker/internal/infrastructure/tracker/generated"
)

type mockIssueQuerier struct {
	mock.Mock
}

func (m *mockIssueQuerier) CreateIssue(ctx context.Context, arg trackerdb.CreateIssueParams) (trackerdb.Issue, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(trackerdb.Issue), args.Error(1)
}

func (m *mockIssueQuerier) ListIssues(ctx context.Context, arg trackerdb.ListIssuesParams) ([]trackerdb.Issue, error) {
	args := m.Called(ctx, arg)
	if result, ok := args.Get(0).([]trackerdb.Issue); ok {
		return result, args.Error(1)
	}
	return nil, args.Error(1)
}

// — CreateIssue unit tests —

func Test_CreateIssue_Success_ReturnsDomainIssue(t *testing.T) {
	querier := &mockIssueQuerier{}
	repository := &IssueRepository{queries: querier}

	projectID := uuid.New()
	reporterID := uuid.New()
	now := time.Now().UTC()

	domainIssue := issuedomain.Issue{
		ID:         uuid.New(),
		Identifier: "test-issue-abc",
		Title:      "Test issue",
		Status:     issuedomain.StatusTodo,
		Priority:   issuedomain.PriorityMedium,
		Labels:     []string{"backend"},
		ProjectID:  projectID,
		ReporterID: reporterID,
	}

	returnedRow := trackerdb.Issue{
		ID:         domainIssue.ID,
		Identifier: domainIssue.Identifier,
		Title:      domainIssue.Title,
		Status:     "todo",
		Priority:   "medium",
		Labels:     []string{"backend"},
		ProjectID:  projectID,
		ReporterID: reporterID,
		CreatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
	}

	querier.On("CreateIssue", mock.Anything, mock.Anything).Return(returnedRow, nil)

	actual, err := repository.CreateIssue(context.Background(), domainIssue)

	require.NoError(t, err)
	require.NotNil(t, actual)
	assert.Equal(t, domainIssue.ID, actual.ID)
	assert.Equal(t, domainIssue.Title, actual.Title)
	assert.Equal(t, issuedomain.StatusTodo, actual.Status)
	querier.AssertExpectations(t)
}

func Test_CreateIssue_QueryError_ReturnsWrappedError(t *testing.T) {
	querier := &mockIssueQuerier{}
	repository := &IssueRepository{queries: querier}

	dbErr := errors.New("unique constraint violation")
	querier.On("CreateIssue", mock.Anything, mock.Anything).Return(trackerdb.Issue{}, dbErr)

	_, err := repository.CreateIssue(context.Background(), issuedomain.Issue{
		ID:         uuid.New(),
		Identifier: "err-test",
		Title:      "Error test",
		Status:     issuedomain.StatusBacklog,
		Priority:   issuedomain.PriorityNone,
		Labels:     []string{},
		ProjectID:  uuid.New(),
		ReporterID: uuid.New(),
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.Contains(t, err.Error(), "create issue")
	querier.AssertExpectations(t)
}

// — ListIssues unit tests —

func Test_ListIssues_Success_ReturnsDomainPage(t *testing.T) {
	querier := &mockIssueQuerier{}
	repository := &IssueRepository{queries: querier}

	projectID := uuid.New()
	now := time.Now().UTC()

	returnedRows := []trackerdb.Issue{
		{
			ID:         uuid.New(),
			Identifier: "issue-1",
			Title:      "First issue",
			Status:     "backlog",
			Priority:   "none",
			Labels:     []string{},
			ProjectID:  projectID,
			ReporterID: uuid.New(),
			CreatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
			UpdatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
		},
	}

	querier.On("ListIssues", mock.Anything, mock.Anything).Return(returnedRows, nil)

	query := issuedomain.ListIssueQuery{}
	actual, err := repository.ListIssues(context.Background(), projectID, query)

	require.NoError(t, err)
	require.Len(t, actual.Items, 1)
	assert.Equal(t, "First issue", actual.Items[0].Title)
	assert.Equal(t, issuedomain.StatusBacklog, actual.Items[0].Status)
	querier.AssertExpectations(t)
}

func Test_ListIssues_EmptyResult_ReturnsEmptyPage(t *testing.T) {
	querier := &mockIssueQuerier{}
	repository := &IssueRepository{queries: querier}

	querier.On("ListIssues", mock.Anything, mock.Anything).Return([]trackerdb.Issue{}, nil)

	query := issuedomain.ListIssueQuery{}
	actual, err := repository.ListIssues(context.Background(), uuid.New(), query)

	require.NoError(t, err)
	assert.Empty(t, actual.Items)
	querier.AssertExpectations(t)
}

func Test_ListIssues_QueryError_ReturnsWrappedError(t *testing.T) {
	querier := &mockIssueQuerier{}
	repository := &IssueRepository{queries: querier}

	dbErr := errors.New("connection refused")
	querier.On("ListIssues", mock.Anything, mock.Anything).Return([]trackerdb.Issue(nil), dbErr)

	query := issuedomain.ListIssueQuery{}
	_, err := repository.ListIssues(context.Background(), uuid.New(), query)

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.Contains(t, err.Error(), "list issues")
	querier.AssertExpectations(t)
}

func Test_ListIssues_WithFilters_PassesCorrectParams(t *testing.T) {
	querier := &mockIssueQuerier{}
	repository := &IssueRepository{queries: querier}

	projectID := uuid.New()
	assigneeID := uuid.New()
	status := issuedomain.StatusInProgress
	priority := issuedomain.PriorityHigh

	expectedParams := listIssuesParamsFromDomain(projectID, issuedomain.ListIssueQuery{
		Status:     &status,
		Priority:   &priority,
		AssigneeID: &assigneeID,
	})

	querier.On("ListIssues", mock.Anything, expectedParams).Return([]trackerdb.Issue{}, nil)

	query := issuedomain.ListIssueQuery{
		Status:     &status,
		Priority:   &priority,
		AssigneeID: &assigneeID,
	}
	_, err := repository.ListIssues(context.Background(), projectID, query)

	require.NoError(t, err)
	querier.AssertExpectations(t)
}
