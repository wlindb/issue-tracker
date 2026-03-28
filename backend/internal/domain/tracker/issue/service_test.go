package issue_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
)

type mockIssueRepository struct {
	mock.Mock
}

func (m *mockIssueRepository) ListIssues(ctx context.Context, projectID uuid.UUID, query issue.ListIssueQuery) (issue.IssuePage, error) {
	args := m.Called(ctx, projectID, query)
	if page, ok := args.Get(0).(issue.IssuePage); ok {
		return page, args.Error(1)
	}
	return issue.IssuePage{}, args.Error(1)
}

func (m *mockIssueRepository) CreateIssue(ctx context.Context, i issue.Issue) (*issue.Issue, error) {
	args := m.Called(ctx, i)
	if result, ok := args.Get(0).(*issue.Issue); ok {
		return result, args.Error(1)
	}
	return nil, args.Error(1)
}

func Test_ListIssues_WithIssues_ReturnsPage(t *testing.T) {
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(repository)

	projectID := uuid.New()
	query := issue.ListIssueQuery{}
	expected := issue.IssuePage{
		Items: []issue.Issue{
			{ID: uuid.New(), Title: "Fix login bug"},
			{ID: uuid.New(), Title: "Add dark mode"},
		},
	}

	repository.On("ListIssues", mock.Anything, projectID, query).Return(expected, nil)

	actual, err := service.ListIssues(context.Background(), projectID, query)
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
	repository.AssertExpectations(t)
}

func Test_ListIssues_EmptyResult_ReturnsEmptyPage(t *testing.T) {
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(repository)

	projectID := uuid.New()
	query := issue.ListIssueQuery{}
	expected := issue.IssuePage{Items: []issue.Issue{}}

	repository.On("ListIssues", mock.Anything, projectID, query).Return(expected, nil)

	actual, err := service.ListIssues(context.Background(), projectID, query)
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
	repository.AssertExpectations(t)
}

func Test_ListIssues_RepositoryError_ReturnsError(t *testing.T) {
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(repository)

	projectID := uuid.New()
	query := issue.ListIssueQuery{}
	repositoryErr := errors.New("db error")

	repository.On("ListIssues", mock.Anything, projectID, query).Return(issue.IssuePage{}, repositoryErr)

	_, err := service.ListIssues(context.Background(), projectID, query)
	require.Error(t, err)
	assert.ErrorIs(t, err, repositoryErr)
	repository.AssertExpectations(t)
}

func Test_ListIssues_WithQueryFilters_PassesFiltersToRepository(t *testing.T) {
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(repository)

	projectID := uuid.New()
	assigneeID := uuid.New()
	status := issue.StatusInProgress
	priority := issue.PriorityHigh
	cursor := "cursor-token"
	limit := 10
	query := issue.ListIssueQuery{
		Cursor:     &cursor,
		Limit:      &limit,
		Status:     &status,
		Priority:   &priority,
		AssigneeID: &assigneeID,
	}
	expected := issue.IssuePage{
		Items: []issue.Issue{
			{ID: uuid.New(), Title: "In-progress issue"},
		},
		NextCursor: nil,
	}

	repository.On("ListIssues", mock.Anything, projectID, query).Return(expected, nil)

	actual, err := service.ListIssues(context.Background(), projectID, query)
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
	repository.AssertExpectations(t)
}

func Test_CreateIssue_ValidCommand_ReturnsCreatedIssue(t *testing.T) {
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(repository)

	projectID := uuid.New()
	reporterID := uuid.New()
	command := issue.CreateIssueCommand{
		ProjectID:  projectID,
		ReporterID: reporterID,
		Title:      "New feature",
		Status:     issue.StatusTodo,
		Priority:   issue.PriorityMedium,
	}
	returned := &issue.Issue{
		ID:         uuid.New(),
		Identifier: "PROJ-1",
		ProjectID:  projectID,
		ReporterID: reporterID,
		Title:      "New feature",
		Status:     issue.StatusTodo,
		Priority:   issue.PriorityMedium,
	}

	repository.On("CreateIssue", mock.Anything, mock.Anything).Return(returned, nil)

	got, err := service.CreateIssue(context.Background(), command)
	require.NoError(t, err)
	assert.Equal(t, returned, got)
	repository.AssertExpectations(t)
}

func Test_CreateIssue_RepositoryError_ReturnsError(t *testing.T) {
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(repository)

	command := issue.CreateIssueCommand{
		ProjectID:  uuid.New(),
		ReporterID: uuid.New(),
		Title:      "New feature",
		Status:     issue.StatusTodo,
		Priority:   issue.PriorityMedium,
	}
	repositoryErr := errors.New("db error")

	repository.On("CreateIssue", mock.Anything, mock.Anything).Return(nil, repositoryErr)

	_, err := service.CreateIssue(context.Background(), command)
	require.Error(t, err)
	assert.ErrorIs(t, err, repositoryErr)
	repository.AssertExpectations(t)
}
