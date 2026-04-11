//go:build !integration

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

func (m *mockIssueRepository) GetIssue(ctx context.Context, id uuid.UUID) (issue.Issue, error) {
	args := m.Called(ctx, id)
	if result, ok := args.Get(0).(issue.Issue); ok {
		return result, args.Error(1)
	}
	return issue.Issue{}, args.Error(1)
}

func (m *mockIssueRepository) ListIssues(ctx context.Context, projectID uuid.UUID, query issue.ListIssueQuery) (issue.IssuePage, error) {
	args := m.Called(ctx, projectID, query)
	if page, ok := args.Get(0).(issue.IssuePage); ok {
		return page, args.Error(1)
	}
	return issue.IssuePage{}, args.Error(1)
}

func (m *mockIssueRepository) CreateIssue(ctx context.Context, i issue.Issue) (issue.Issue, error) {
	args := m.Called(ctx, i)
	if result, ok := args.Get(0).(issue.Issue); ok {
		return result, args.Error(1)
	}
	return issue.Issue{}, args.Error(1)
}

func (m *mockIssueRepository) Update(ctx context.Context, i issue.Issue) (issue.Issue, error) {
	args := m.Called(ctx, i)
	if result, ok := args.Get(0).(issue.Issue); ok {
		return result, args.Error(1)
	}
	return issue.Issue{}, args.Error(1)
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
	returned := issue.Issue{
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

// — GetIssue —

func Test_GetIssue_Found_ReturnsIssue(t *testing.T) {
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(repository)

	issueID := uuid.New()
	expected := issue.Issue{
		ID:       issueID,
		Title:    "Test issue",
		Status:   issue.StatusTodo,
		Priority: issue.PriorityNone,
	}

	repository.On("GetIssue", mock.Anything, issueID).Return(expected, nil)

	actual, err := service.GetIssue(context.Background(), issueID)
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
	repository.AssertExpectations(t)
}

func Test_GetIssue_NotFound_ReturnsError(t *testing.T) {
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(repository)

	issueID := uuid.New()
	repository.On("GetIssue", mock.Anything, issueID).Return(issue.Issue{}, issue.ErrIssueNotFound)

	_, err := service.GetIssue(context.Background(), issueID)
	require.Error(t, err)
	assert.ErrorIs(t, err, issue.ErrIssueNotFound)
	repository.AssertExpectations(t)
}

// — UpdateIssueAssignee —

func Test_UpdateIssueAssignee_ValidAssignee_ReturnsUpdatedIssue(t *testing.T) {
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(repository)

	issueID := uuid.New()
	assigneeID := uuid.New()
	existing := issue.Issue{
		ID:       issueID,
		Title:    "Test issue",
		Status:   issue.StatusTodo,
		Priority: issue.PriorityNone,
		Labels:   []string{},
	}
	returned := existing
	returned.AssigneeID = &assigneeID

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)
	repository.On("Update", mock.Anything, mock.MatchedBy(func(i issue.Issue) bool {
		return i.ID == issueID && i.AssigneeID != nil && *i.AssigneeID == assigneeID
	})).Return(returned, nil)

	actual, err := service.UpdateIssueAssignee(context.Background(), issueID, &assigneeID)
	require.NoError(t, err)
	require.NotNil(t, actual.AssigneeID)
	assert.Equal(t, assigneeID, *actual.AssigneeID)
	repository.AssertExpectations(t)
}

func Test_UpdateIssueAssignee_NilAssignee_ClearsAssignee(t *testing.T) {
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(repository)

	issueID := uuid.New()
	oldAssignee := uuid.New()
	existing := issue.Issue{
		ID:         issueID,
		Title:      "Test issue",
		Status:     issue.StatusTodo,
		Priority:   issue.PriorityNone,
		Labels:     []string{},
		AssigneeID: &oldAssignee,
	}
	returned := existing
	returned.AssigneeID = nil

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)
	repository.On("Update", mock.Anything, mock.MatchedBy(func(i issue.Issue) bool {
		return i.ID == issueID && i.AssigneeID == nil
	})).Return(returned, nil)

	actual, err := service.UpdateIssueAssignee(context.Background(), issueID, nil)
	require.NoError(t, err)
	assert.Nil(t, actual.AssigneeID)
	repository.AssertExpectations(t)
}

func Test_UpdateIssueAssignee_IssueNotFound_ReturnsError(t *testing.T) {
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(repository)

	issueID := uuid.New()
	assigneeID := uuid.New()
	repository.On("GetIssue", mock.Anything, issueID).Return(issue.Issue{}, issue.ErrIssueNotFound)

	_, err := service.UpdateIssueAssignee(context.Background(), issueID, &assigneeID)
	require.Error(t, err)
	assert.ErrorIs(t, err, issue.ErrIssueNotFound)
	repository.AssertExpectations(t)
}

func Test_UpdateIssueAssignee_UpdateError_ReturnsError(t *testing.T) {
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(repository)

	issueID := uuid.New()
	assigneeID := uuid.New()
	existing := issue.Issue{ID: issueID, Status: issue.StatusTodo, Priority: issue.PriorityNone, Labels: []string{}}
	updateErr := errors.New("db error")

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)
	repository.On("Update", mock.Anything, mock.Anything).Return(issue.Issue{}, updateErr)

	_, err := service.UpdateIssueAssignee(context.Background(), issueID, &assigneeID)
	require.Error(t, err)
	assert.ErrorIs(t, err, updateErr)
	repository.AssertExpectations(t)
}

// — UpdateIssueDescription —

func Test_UpdateIssueDescription_ValidDescription_ReturnsUpdatedIssue(t *testing.T) {
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(repository)

	issueID := uuid.New()
	description := "new description"
	existing := issue.Issue{
		ID:       issueID,
		Title:    "Test issue",
		Status:   issue.StatusTodo,
		Priority: issue.PriorityNone,
		Labels:   []string{},
	}
	returned := existing
	returned.Description = &description

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)
	repository.On("Update", mock.Anything, mock.MatchedBy(func(i issue.Issue) bool {
		return i.ID == issueID && i.Description != nil && *i.Description == description
	})).Return(returned, nil)

	actual, err := service.UpdateIssueDescription(context.Background(), issueID, &description)
	require.NoError(t, err)
	require.NotNil(t, actual.Description)
	assert.Equal(t, description, *actual.Description)
	repository.AssertExpectations(t)
}

func Test_UpdateIssueDescription_NilDescription_ClearsDescription(t *testing.T) {
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(repository)

	issueID := uuid.New()
	oldDesc := "old description"
	existing := issue.Issue{
		ID:          issueID,
		Title:       "Test issue",
		Description: &oldDesc,
		Status:      issue.StatusTodo,
		Priority:    issue.PriorityNone,
		Labels:      []string{},
	}
	returned := existing
	returned.Description = nil

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)
	repository.On("Update", mock.Anything, mock.MatchedBy(func(i issue.Issue) bool {
		return i.ID == issueID && i.Description == nil
	})).Return(returned, nil)

	actual, err := service.UpdateIssueDescription(context.Background(), issueID, nil)
	require.NoError(t, err)
	assert.Nil(t, actual.Description)
	repository.AssertExpectations(t)
}

func Test_UpdateIssueDescription_IssueNotFound_ReturnsError(t *testing.T) {
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(repository)

	issueID := uuid.New()
	desc := "description"
	repository.On("GetIssue", mock.Anything, issueID).Return(issue.Issue{}, issue.ErrIssueNotFound)

	_, err := service.UpdateIssueDescription(context.Background(), issueID, &desc)
	require.Error(t, err)
	assert.ErrorIs(t, err, issue.ErrIssueNotFound)
	repository.AssertExpectations(t)
}

// — UpdateIssuePriority —

func Test_UpdateIssuePriority_ValidPriority_ReturnsUpdatedIssue(t *testing.T) {
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(repository)

	issueID := uuid.New()
	existing := issue.Issue{
		ID:       issueID,
		Title:    "Test issue",
		Status:   issue.StatusTodo,
		Priority: issue.PriorityNone,
		Labels:   []string{},
	}
	returned := existing
	returned.Priority = issue.PriorityHigh

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)
	repository.On("Update", mock.Anything, mock.MatchedBy(func(i issue.Issue) bool {
		return i.ID == issueID && i.Priority == issue.PriorityHigh
	})).Return(returned, nil)

	actual, err := service.UpdateIssuePriority(context.Background(), issueID, issue.PriorityHigh)
	require.NoError(t, err)
	assert.Equal(t, issue.PriorityHigh, actual.Priority)
	repository.AssertExpectations(t)
}

func Test_UpdateIssuePriority_InvalidPriority_ReturnsError(t *testing.T) {
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(repository)

	issueID := uuid.New()
	existing := issue.Issue{
		ID:       issueID,
		Title:    "Test issue",
		Status:   issue.StatusTodo,
		Priority: issue.PriorityNone,
		Labels:   []string{},
	}

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)

	_, err := service.UpdateIssuePriority(context.Background(), issueID, issue.Priority("extreme"))
	require.Error(t, err)
	assert.ErrorIs(t, err, issue.ErrInvalidIssue)
	repository.AssertExpectations(t)
}

func Test_UpdateIssuePriority_IssueNotFound_ReturnsError(t *testing.T) {
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(repository)

	issueID := uuid.New()
	repository.On("GetIssue", mock.Anything, issueID).Return(issue.Issue{}, issue.ErrIssueNotFound)

	_, err := service.UpdateIssuePriority(context.Background(), issueID, issue.PriorityHigh)
	require.Error(t, err)
	assert.ErrorIs(t, err, issue.ErrIssueNotFound)
	repository.AssertExpectations(t)
}

func Test_UpdateIssuePriority_UpdateError_ReturnsError(t *testing.T) {
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(repository)

	issueID := uuid.New()
	existing := issue.Issue{ID: issueID, Status: issue.StatusTodo, Priority: issue.PriorityNone, Labels: []string{}}
	updateErr := errors.New("db error")

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)
	repository.On("Update", mock.Anything, mock.Anything).Return(issue.Issue{}, updateErr)

	_, err := service.UpdateIssuePriority(context.Background(), issueID, issue.PriorityHigh)
	require.Error(t, err)
	assert.ErrorIs(t, err, updateErr)
	repository.AssertExpectations(t)
}

// — UpdateIssueStatus —

func Test_UpdateIssueStatus_ValidStatus_ReturnsUpdatedIssue(t *testing.T) {
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(repository)

	issueID := uuid.New()
	existing := issue.Issue{
		ID:       issueID,
		Title:    "Test issue",
		Status:   issue.StatusTodo,
		Priority: issue.PriorityNone,
		Labels:   []string{},
	}
	returned := existing
	returned.Status = issue.StatusDone

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)
	repository.On("Update", mock.Anything, mock.MatchedBy(func(i issue.Issue) bool {
		return i.ID == issueID && i.Status == issue.StatusDone
	})).Return(returned, nil)

	actual, err := service.UpdateIssueStatus(context.Background(), issueID, issue.StatusDone)
	require.NoError(t, err)
	assert.Equal(t, issue.StatusDone, actual.Status)
	repository.AssertExpectations(t)
}

func Test_UpdateIssueStatus_InvalidStatus_ReturnsError(t *testing.T) {
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(repository)

	issueID := uuid.New()
	existing := issue.Issue{
		ID:       issueID,
		Title:    "Test issue",
		Status:   issue.StatusTodo,
		Priority: issue.PriorityNone,
		Labels:   []string{},
	}

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)

	_, err := service.UpdateIssueStatus(context.Background(), issueID, issue.Status("archived"))
	require.Error(t, err)
	assert.ErrorIs(t, err, issue.ErrInvalidIssue)
	repository.AssertExpectations(t)
}

func Test_UpdateIssueStatus_IssueNotFound_ReturnsError(t *testing.T) {
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(repository)

	issueID := uuid.New()
	repository.On("GetIssue", mock.Anything, issueID).Return(issue.Issue{}, issue.ErrIssueNotFound)

	_, err := service.UpdateIssueStatus(context.Background(), issueID, issue.StatusDone)
	require.Error(t, err)
	assert.ErrorIs(t, err, issue.ErrIssueNotFound)
	repository.AssertExpectations(t)
}

func Test_UpdateIssueStatus_UpdateError_ReturnsError(t *testing.T) {
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(repository)

	issueID := uuid.New()
	existing := issue.Issue{ID: issueID, Status: issue.StatusTodo, Priority: issue.PriorityNone, Labels: []string{}}
	updateErr := errors.New("db error")

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)
	repository.On("Update", mock.Anything, mock.Anything).Return(issue.Issue{}, updateErr)

	_, err := service.UpdateIssueStatus(context.Background(), issueID, issue.StatusDone)
	require.Error(t, err)
	assert.ErrorIs(t, err, updateErr)
	repository.AssertExpectations(t)
}
