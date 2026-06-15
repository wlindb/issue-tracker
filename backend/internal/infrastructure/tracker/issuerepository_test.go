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

func (m *mockIssueQuerier) GetIssue(ctx context.Context, id uuid.UUID) (trackerdb.Issue, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(trackerdb.Issue), args.Error(1)
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

func (m *mockIssueQuerier) UpdateIssue(ctx context.Context, arg trackerdb.UpdateIssueParams) (trackerdb.Issue, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(trackerdb.Issue), args.Error(1)
}

func (m *mockIssueQuerier) ListLabelsByIssueIDs(ctx context.Context, issueIDs []uuid.UUID) ([]trackerdb.ListLabelsByIssueIDsRow, error) {
	args := m.Called(ctx, issueIDs)
	if result, ok := args.Get(0).([]trackerdb.ListLabelsByIssueIDsRow); ok {
		return result, args.Error(1)
	}
	return []trackerdb.ListLabelsByIssueIDsRow{}, args.Error(1)
}

func (m *mockIssueQuerier) GetLabelsByIDs(ctx context.Context, ids []uuid.UUID) ([]trackerdb.GetLabelsByIDsRow, error) {
	args := m.Called(ctx, ids)
	if result, ok := args.Get(0).([]trackerdb.GetLabelsByIDsRow); ok {
		return result, args.Error(1)
	}
	return []trackerdb.GetLabelsByIDsRow{}, args.Error(1)
}

// — ListIssues unit tests —

func Test_ListIssues_Success_ReturnsDomainPage(t *testing.T) {
	querier := &mockIssueQuerier{}
	repository := &IssueRepository{queries: querier}

	projectID := uuid.New()
	issueID := uuid.New()
	labelID := uuid.New()
	now := time.Now().UTC()

	returnedRows := []trackerdb.Issue{
		{
			ID:         issueID,
			Identifier: "issue-1",
			Title:      "First issue",
			Status:     "backlog",
			Priority:   "none",
			ProjectID:  projectID,
			ReporterID: uuid.New(),
			CreatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
			UpdatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
		},
	}
	returnedLabelRows := []trackerdb.ListLabelsByIssueIDsRow{
		{IssueID: issueID, ID: labelID, Name: "backend"},
	}

	querier.On("ListIssues", mock.Anything, mock.Anything).Return(returnedRows, nil)
	querier.On("ListLabelsByIssueIDs", mock.Anything, []uuid.UUID{issueID}).Return(returnedLabelRows, nil)

	query := issuedomain.ListIssueQuery{}
	actual, err := repository.ListIssues(context.Background(), projectID, query)

	require.NoError(t, err)
	require.Len(t, actual.Items, 1)
	assert.Equal(t, "First issue", actual.Items[0].Title)
	assert.Equal(t, issuedomain.StatusBacklog, actual.Items[0].Status)
	require.Len(t, actual.Items[0].Labels, 1)
	assert.Equal(t, labelID, actual.Items[0].Labels[0].ID)
	assert.Equal(t, "backend", actual.Items[0].Labels[0].Name)
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

// — GetIssue unit tests —

func Test_GetIssue_Success_ReturnsDomainIssue(t *testing.T) {
	querier := &mockIssueQuerier{}
	repository := &IssueRepository{queries: querier}

	issueID := uuid.New()
	labelID := uuid.New()
	now := time.Now().UTC()

	returnedRow := trackerdb.Issue{
		ID:         issueID,
		Identifier: "issue-1",
		Title:      "Test issue",
		Status:     "todo",
		Priority:   "medium",
		ProjectID:  uuid.New(),
		ReporterID: uuid.New(),
		CreatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
	}
	returnedLabelRows := []trackerdb.ListLabelsByIssueIDsRow{
		{IssueID: issueID, ID: labelID, Name: "frontend"},
	}

	querier.On("GetIssue", mock.Anything, issueID).Return(returnedRow, nil)
	querier.On("ListLabelsByIssueIDs", mock.Anything, []uuid.UUID{issueID}).Return(returnedLabelRows, nil)

	actual, err := repository.GetIssue(context.Background(), issueID)

	require.NoError(t, err)
	assert.Equal(t, issueID, actual.ID)
	assert.Equal(t, "Test issue", actual.Title)
	require.Len(t, actual.Labels, 1)
	assert.Equal(t, labelID, actual.Labels[0].ID)
	assert.Equal(t, "frontend", actual.Labels[0].Name)
	querier.AssertExpectations(t)
}

// — Update unit tests —

func Test_Update_Success_ReturnsDomainIssue(t *testing.T) {
	querier := &mockIssueQuerier{}
	repository := &IssueRepository{queries: querier}

	projectID := uuid.New()
	reporterID := uuid.New()
	description := "updated desc"
	now := time.Now().UTC()
	labelID := uuid.New()

	domainIssue := issuedomain.Issue{
		ID:          uuid.New(),
		Identifier:  "test-issue-abc",
		Title:       "Test issue",
		Description: &description,
		Status:      issuedomain.StatusInProgress,
		Priority:    issuedomain.PriorityHigh,
		Labels:      []issuedomain.Label{{ID: labelID, Name: "backend"}},
		ProjectID:   projectID,
		ReporterID:  reporterID,
	}

	returnedRow := trackerdb.Issue{
		ID:          domainIssue.ID,
		Identifier:  domainIssue.Identifier,
		Title:       domainIssue.Title,
		Description: pgtype.Text{String: description, Valid: true},
		Status:      "in_progress",
		Priority:    "high",
		ProjectID:   projectID,
		ReporterID:  reporterID,
		CreatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
	}

	querier.On("UpdateIssue", mock.Anything, mock.Anything).Return(returnedRow, nil)

	actual, err := repository.Update(context.Background(), domainIssue)

	require.NoError(t, err)
	assert.Equal(t, domainIssue.ID, actual.ID)
	assert.Equal(t, issuedomain.StatusInProgress, actual.Status)
	assert.Equal(t, issuedomain.PriorityHigh, actual.Priority)
	require.NotNil(t, actual.Description)
	assert.Equal(t, description, *actual.Description)
	require.Len(t, actual.Labels, 1)
	assert.Equal(t, labelID, actual.Labels[0].ID)
	querier.AssertExpectations(t)
}

func Test_Update_QueryError_ReturnsWrappedError(t *testing.T) {
	querier := &mockIssueQuerier{}
	repository := &IssueRepository{queries: querier}

	dbErr := errors.New("update conflict")
	querier.On("UpdateIssue", mock.Anything, mock.Anything).Return(trackerdb.Issue{}, dbErr)

	_, err := repository.Update(context.Background(), issuedomain.Issue{
		ID:       uuid.New(),
		Status:   issuedomain.StatusDone,
		Priority: issuedomain.PriorityNone,
		Labels:   []issuedomain.Label{},
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.Contains(t, err.Error(), "update issue")
	querier.AssertExpectations(t)
}

// — GetLabelsByIDs unit tests —

func Test_GetLabelsByIDs_EmptyIDs_ReturnsEmptySlice(t *testing.T) {
	querier := &mockIssueQuerier{}
	repository := &IssueRepository{queries: querier}

	actual, err := repository.GetLabelsByIDs(context.Background(), []uuid.UUID{})

	require.NoError(t, err)
	assert.Empty(t, actual)
	querier.AssertNotCalled(t, "GetLabelsByIDs")
}

func Test_GetLabelsByIDs_ValidIDs_ReturnsDomainLabels(t *testing.T) {
	querier := &mockIssueQuerier{}
	repository := &IssueRepository{queries: querier}

	id1, id2 := uuid.New(), uuid.New()
	returnedRows := []trackerdb.GetLabelsByIDsRow{
		{ID: id1, Name: "bug"},
		{ID: id2, Name: "feature"},
	}

	querier.On("GetLabelsByIDs", mock.Anything, []uuid.UUID{id1, id2}).Return(returnedRows, nil)

	actual, err := repository.GetLabelsByIDs(context.Background(), []uuid.UUID{id1, id2})

	require.NoError(t, err)
	require.Len(t, actual, 2)
	assert.Equal(t, id1, actual[0].ID)
	assert.Equal(t, "bug", actual[0].Name)
	assert.Equal(t, id2, actual[1].ID)
	assert.Equal(t, "feature", actual[1].Name)
	querier.AssertExpectations(t)
}
