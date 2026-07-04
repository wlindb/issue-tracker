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
	"github.com/wlindb/issue-tracker/internal/domain/tracker/label"
	"github.com/wlindb/issue-tracker/internal/pkg/domain/event"
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

type mockLabelRepository struct {
	mock.Mock
}

func (m *mockLabelRepository) GetOrCreate(ctx context.Context, name string) (label.Label, error) {
	args := m.Called(ctx, name)
	if result, ok := args.Get(0).(label.Label); ok {
		return result, args.Error(1)
	}
	return label.Label{}, args.Error(1)
}

func (m *mockLabelRepository) ListByIDs(ctx context.Context, ids []uuid.UUID) ([]label.Label, error) {
	args := m.Called(ctx, ids)
	if result, ok := args.Get(0).([]label.Label); ok {
		return result, args.Error(1)
	}
	return []label.Label{}, args.Error(1)
}

func (m *mockLabelRepository) SearchByName(ctx context.Context, name string) ([]label.Label, error) {
	args := m.Called(ctx, name)
	if result, ok := args.Get(0).([]label.Label); ok {
		return result, args.Error(1)
	}
	return []label.Label{}, args.Error(1)
}

type mockUnitOfWork struct {
	mock.Mock
}

func (m *mockUnitOfWork) RunInTx(ctx context.Context, fn func(issue.Repositories) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

type fakeUnitOfWork struct {
	repositories issue.Repositories
}

func (f *fakeUnitOfWork) RunInTx(_ context.Context, fn func(issue.Repositories) error) error {
	return fn(f.repositories)
}

func Test_ListIssues_WithIssues_ReturnsPage(t *testing.T) {
	uow := &mockUnitOfWork{}
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

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
	uow := &mockUnitOfWork{}
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

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
	uow := &mockUnitOfWork{}
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

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
	uow := &mockUnitOfWork{}
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

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
	labelRepository := &mockLabelRepository{}
	uow := &fakeUnitOfWork{repositories: issue.Repositories{Issues: repository}}
	service := issue.NewIssueService(uow, repository, labelRepository)

	projectID := uuid.New()
	reporterID := uuid.New()
	command := issue.CreateIssueCommand{
		ProjectID:  projectID,
		ReporterID: reporterID,
		Title:      "New feature",
		Status:     issue.StatusTodo,
		Priority:   issue.PriorityMedium,
	}
	expected := issue.Issue{
		ID:         uuid.New(),
		Identifier: "PROJ-1",
		ProjectID:  projectID,
		ReporterID: reporterID,
		Title:      "New feature",
		Status:     issue.StatusTodo,
		Priority:   issue.PriorityMedium,
	}

	labelRepository.On("ListByIDs", mock.Anything, command.LabelIDs).Return([]label.Label{}, nil)
	repository.On("CreateIssue", mock.Anything, mock.Anything).Return(expected, nil)

	ctxWithNoopPublisher := event.WithPublisher(context.Background(), func(_ context.Context, _ issue.IssueCreatedEvent) error {
		return nil
	})

	actual, err := service.CreateIssue(ctxWithNoopPublisher, command)
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
	labelRepository.AssertExpectations(t)
	repository.AssertExpectations(t)
}

func Test_CreateIssue_RunInTxError_ReturnsError(t *testing.T) {
	repository := &mockIssueRepository{}
	labelRepository := &mockLabelRepository{}
	repositoryErr := errors.New("something went wrong")

	uow := &mockUnitOfWork{}
	uow.On("RunInTx", mock.Anything, mock.Anything).Return(repositoryErr)

	service := issue.NewIssueService(uow, repository, labelRepository)

	command := issue.CreateIssueCommand{
		ProjectID:  uuid.New(),
		ReporterID: uuid.New(),
		Title:      "New feature",
		Status:     issue.StatusTodo,
		Priority:   issue.PriorityMedium,
	}
	labelRepository.On("ListByIDs", mock.Anything, command.LabelIDs).Return([]label.Label{}, nil)

	_, err := service.CreateIssue(context.Background(), command)
	require.Error(t, err)
	assert.ErrorIs(t, err, repositoryErr)
	labelRepository.AssertExpectations(t)
	repository.AssertExpectations(t)
}

func Test_CreateIssue_RepositoryError_ReturnsError(t *testing.T) {
	repositoryErr := errors.New("db error")
	txRepository := &mockIssueRepository{}
	labelRepository := &mockLabelRepository{}
	txRepository.On("CreateIssue", mock.Anything, mock.Anything).Return(issue.Issue{}, repositoryErr)

	uow := &fakeUnitOfWork{repositories: issue.Repositories{Issues: txRepository}}
	service := issue.NewIssueService(uow, txRepository, labelRepository)

	command := issue.CreateIssueCommand{
		ProjectID:  uuid.New(),
		ReporterID: uuid.New(),
		Title:      "New feature",
		Status:     issue.StatusTodo,
		Priority:   issue.PriorityMedium,
	}
	labelRepository.On("ListByIDs", mock.Anything, command.LabelIDs).Return([]label.Label{}, nil)

	_, err := service.CreateIssue(context.Background(), command)
	require.Error(t, err)
	assert.ErrorIs(t, err, repositoryErr)
	labelRepository.AssertExpectations(t)
	txRepository.AssertExpectations(t)
}

func Test_CreateIssue_WithLabelIDs_ResolvesLabelsBeforeTransaction(t *testing.T) {
	repository := &mockIssueRepository{}
	labelRepository := &mockLabelRepository{}
	uow := &fakeUnitOfWork{repositories: issue.Repositories{Issues: repository}}
	service := issue.NewIssueService(uow, repository, labelRepository)

	projectID := uuid.New()
	reporterID := uuid.New()
	labelIDOne := uuid.New()
	labelIDTwo := uuid.New()
	command := issue.CreateIssueCommand{
		ProjectID:  projectID,
		ReporterID: reporterID,
		Title:      "New feature",
		Status:     issue.StatusTodo,
		Priority:   issue.PriorityMedium,
		LabelIDs:   []uuid.UUID{labelIDOne, labelIDTwo},
	}
	resolvedLabels := []label.Label{
		{ID: labelIDOne, Name: "bug"},
		{ID: labelIDTwo, Name: "high-priority"},
	}
	expected := issue.Issue{
		ID:         uuid.New(),
		Identifier: "PROJ-1",
		ProjectID:  projectID,
		ReporterID: reporterID,
		Title:      "New feature",
		Status:     issue.StatusTodo,
		Priority:   issue.PriorityMedium,
		Labels:     resolvedLabels,
	}

	labelRepository.On("ListByIDs", mock.Anything, command.LabelIDs).Return(resolvedLabels, nil)
	repository.On("CreateIssue", mock.Anything, mock.Anything).Return(expected, nil)

	ctxWithNoopPublisher := event.WithPublisher(context.Background(), func(_ context.Context, _ issue.IssueCreatedEvent) error {
		return nil
	})

	actual, err := service.CreateIssue(ctxWithNoopPublisher, command)
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
	labelRepository.AssertExpectations(t)
	repository.AssertExpectations(t)
}

func Test_CreateIssue_LabelRepositoryError_ReturnsError(t *testing.T) {
	repository := &mockIssueRepository{}
	labelRepository := &mockLabelRepository{}
	uow := &mockUnitOfWork{}
	service := issue.NewIssueService(uow, repository, labelRepository)

	command := issue.CreateIssueCommand{
		ProjectID:  uuid.New(),
		ReporterID: uuid.New(),
		Title:      "New feature",
		Status:     issue.StatusTodo,
		Priority:   issue.PriorityMedium,
		LabelIDs:   []uuid.UUID{uuid.New()},
	}
	labelRepositoryErr := errors.New("list labels failed")
	labelRepository.On("ListByIDs", mock.Anything, command.LabelIDs).Return([]label.Label{}, labelRepositoryErr)

	_, err := service.CreateIssue(context.Background(), command)
	require.Error(t, err)
	assert.ErrorIs(t, err, labelRepositoryErr)
	uow.AssertNotCalled(t, "RunInTx", mock.Anything, mock.Anything)
	labelRepository.AssertExpectations(t)
	repository.AssertNotCalled(t, "CreateIssue", mock.Anything, mock.Anything)
}

// — GetIssue —

func Test_GetIssue_Found_ReturnsIssue(t *testing.T) {
	uow := &mockUnitOfWork{}
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

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
	uow := &mockUnitOfWork{}
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

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
	uow := &fakeUnitOfWork{repositories: issue.Repositories{Issues: repository}}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

	issueID := uuid.New()
	assigneeID := uuid.New()
	existing := issue.Issue{
		ID:       issueID,
		Title:    "Test issue",
		Status:   issue.StatusTodo,
		Priority: issue.PriorityNone,
		Labels:   []label.Label{},
	}
	returned := existing
	returned.AssigneeID = &assigneeID

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)
	repository.On("Update", mock.Anything, mock.MatchedBy(func(i issue.Issue) bool {
		return i.ID == issueID && i.AssigneeID != nil && *i.AssigneeID == assigneeID
	})).Return(returned, nil)

	ctxWithNoopPublisher := event.WithPublisher(context.Background(), func(_ context.Context, _ issue.IssueAssigneeUpdatedEvent) error {
		return nil
	})

	actual, err := service.UpdateIssueAssignee(ctxWithNoopPublisher, issueID, &assigneeID)
	require.NoError(t, err)
	require.NotNil(t, actual.AssigneeID)
	assert.Equal(t, assigneeID, *actual.AssigneeID)
	repository.AssertExpectations(t)
}

func Test_UpdateIssueAssignee_NilAssignee_ClearsAssignee(t *testing.T) {
	repository := &mockIssueRepository{}
	uow := &fakeUnitOfWork{repositories: issue.Repositories{Issues: repository}}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

	issueID := uuid.New()
	oldAssignee := uuid.New()
	existing := issue.Issue{
		ID:         issueID,
		Title:      "Test issue",
		Status:     issue.StatusTodo,
		Priority:   issue.PriorityNone,
		Labels:     []label.Label{},
		AssigneeID: &oldAssignee,
	}
	returned := existing
	returned.AssigneeID = nil

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)
	repository.On("Update", mock.Anything, mock.MatchedBy(func(i issue.Issue) bool {
		return i.ID == issueID && i.AssigneeID == nil
	})).Return(returned, nil)

	ctxWithNoopPublisher := event.WithPublisher(context.Background(), func(_ context.Context, _ issue.IssueAssigneeUpdatedEvent) error {
		return nil
	})

	actual, err := service.UpdateIssueAssignee(ctxWithNoopPublisher, issueID, nil)
	require.NoError(t, err)
	assert.Nil(t, actual.AssigneeID)
	repository.AssertExpectations(t)
}

func Test_UpdateIssueAssignee_IssueNotFound_ReturnsError(t *testing.T) {
	repository := &mockIssueRepository{}
	uow := &fakeUnitOfWork{repositories: issue.Repositories{Issues: repository}}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

	issueID := uuid.New()
	assigneeID := uuid.New()
	repository.On("GetIssue", mock.Anything, issueID).Return(issue.Issue{}, issue.ErrIssueNotFound)

	ctxWithNoopPublisher := event.WithPublisher(context.Background(), func(_ context.Context, _ issue.IssueAssigneeUpdatedEvent) error {
		return nil
	})

	_, err := service.UpdateIssueAssignee(ctxWithNoopPublisher, issueID, &assigneeID)
	require.Error(t, err)
	assert.ErrorIs(t, err, issue.ErrIssueNotFound)
	repository.AssertExpectations(t)
}

func Test_UpdateIssueAssignee_UpdateError_ReturnsError(t *testing.T) {
	repository := &mockIssueRepository{}
	uow := &fakeUnitOfWork{repositories: issue.Repositories{Issues: repository}}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

	issueID := uuid.New()
	assigneeID := uuid.New()
	existing := issue.Issue{ID: issueID, Status: issue.StatusTodo, Priority: issue.PriorityNone, Labels: []label.Label{}}
	updateErr := errors.New("db error")

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)
	repository.On("Update", mock.Anything, mock.Anything).Return(issue.Issue{}, updateErr)

	_, err := service.UpdateIssueAssignee(context.Background(), issueID, &assigneeID)
	require.Error(t, err)
	assert.ErrorIs(t, err, updateErr)
	repository.AssertExpectations(t)
}

func Test_UpdateIssueAssignee_SuccessfulUpdate_PublishesIssueAssigneeUpdatedEvent(t *testing.T) {
	repository := &mockIssueRepository{}
	uow := &fakeUnitOfWork{repositories: issue.Repositories{Issues: repository}}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

	issueID := uuid.New()
	assigneeID := uuid.New()
	existing := issue.Issue{
		ID:       issueID,
		Title:    "Test issue",
		Status:   issue.StatusTodo,
		Priority: issue.PriorityNone,
		Labels:   []label.Label{},
	}
	returned := existing
	returned.AssigneeID = &assigneeID

	var published []issue.IssueAssigneeUpdatedEvent
	ctx := event.WithPublisher(context.Background(), func(_ context.Context, e issue.IssueAssigneeUpdatedEvent) error {
		published = append(published, e)
		return nil
	})

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)
	repository.On("Update", mock.Anything, mock.MatchedBy(func(i issue.Issue) bool {
		return i.ID == issueID && i.AssigneeID != nil && *i.AssigneeID == assigneeID
	})).Return(returned, nil)

	actual, err := service.UpdateIssueAssignee(ctx, issueID, &assigneeID)
	require.NoError(t, err)
	require.Len(t, published, 1)
	assert.Equal(t, returned, actual)
	assert.Equal(t, returned, published[0].Payload)
	repository.AssertExpectations(t)
}

func Test_UpdateIssueAssignee_EmitAssigneeUpdatedError_ReturnsError(t *testing.T) {
	repository := &mockIssueRepository{}
	uow := &fakeUnitOfWork{repositories: issue.Repositories{Issues: repository}}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

	issueID := uuid.New()
	assigneeID := uuid.New()
	existing := issue.Issue{
		ID:       issueID,
		Title:    "Test issue",
		Status:   issue.StatusTodo,
		Priority: issue.PriorityNone,
		Labels:   []label.Label{},
	}
	updated := existing
	updated.AssigneeID = &assigneeID

	expectedError := errors.New("publish error")
	ctx := event.WithPublisher(context.Background(), func(_ context.Context, _ issue.IssueAssigneeUpdatedEvent) error {
		return expectedError
	})

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)
	repository.On("Update", mock.Anything, mock.Anything).Return(updated, nil)

	_, err := service.UpdateIssueAssignee(ctx, issueID, &assigneeID)
	assert.ErrorIs(t, err, expectedError)
	repository.AssertExpectations(t)
}

// — UpdateIssueDescription —

func Test_UpdateIssueDescription_ValidDescription_ReturnsUpdatedIssue(t *testing.T) {
	repository := &mockIssueRepository{}
	uow := &fakeUnitOfWork{repositories: issue.Repositories{Issues: repository}}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

	issueID := uuid.New()
	description := "new description"
	existing := issue.Issue{
		ID:       issueID,
		Title:    "Test issue",
		Status:   issue.StatusTodo,
		Priority: issue.PriorityNone,
		Labels:   []label.Label{},
	}
	returned := existing
	returned.Description = &description

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)
	repository.On("Update", mock.Anything, mock.Anything).Return(returned, nil)

	ctxWithNoopPublisher := event.WithPublisher(context.Background(), func(_ context.Context, _ issue.IssueDescriptionUpdatedEvent) error {
		return nil
	})

	actual, err := service.UpdateIssueDescription(ctxWithNoopPublisher, issueID, &description)
	require.NoError(t, err)
	require.NotNil(t, actual.Description)
	assert.Equal(t, description, *actual.Description)
	repository.AssertExpectations(t)
}

func Test_UpdateIssueDescription_NilDescription_ClearsDescription(t *testing.T) {
	repository := &mockIssueRepository{}
	uow := &fakeUnitOfWork{repositories: issue.Repositories{Issues: repository}}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

	issueID := uuid.New()
	oldDesc := "old description"
	existing := issue.Issue{
		ID:          issueID,
		Title:       "Test issue",
		Description: &oldDesc,
		Status:      issue.StatusTodo,
		Priority:    issue.PriorityNone,
		Labels:      []label.Label{},
	}
	returned := existing
	returned.Description = nil

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)
	repository.On("Update", mock.Anything, mock.MatchedBy(func(i issue.Issue) bool {
		return i.ID == issueID && i.Description == nil
	})).Return(returned, nil)

	ctxWithNoopPublisher := event.WithPublisher(context.Background(), func(_ context.Context, _ issue.IssueDescriptionUpdatedEvent) error {
		return nil
	})

	actual, err := service.UpdateIssueDescription(ctxWithNoopPublisher, issueID, nil)
	require.NoError(t, err)
	assert.Nil(t, actual.Description)
	repository.AssertExpectations(t)
}

func Test_UpdateIssueDescription_IssueNotFound_ReturnsError(t *testing.T) {
	repository := &mockIssueRepository{}
	uow := &fakeUnitOfWork{repositories: issue.Repositories{Issues: repository}}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

	issueID := uuid.New()
	desc := "description"
	repository.On("GetIssue", mock.Anything, issueID).Return(issue.Issue{}, issue.ErrIssueNotFound)

	_, err := service.UpdateIssueDescription(context.Background(), issueID, &desc)
	require.Error(t, err)
	assert.ErrorIs(t, err, issue.ErrIssueNotFound)
	repository.AssertExpectations(t)
}

func Test_UpdateIssueDescription_SuccessfulUpdate_PublishesIssueDescriptionUpdatedEvent(t *testing.T) {
	repository := &mockIssueRepository{}
	uow := &fakeUnitOfWork{repositories: issue.Repositories{Issues: repository}}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

	issueID := uuid.New()
	oldDesc := "old description"
	newDesc := "new description"
	existing := issue.Issue{
		ID:          issueID,
		Title:       "Test issue",
		Description: &oldDesc,
		Status:      issue.StatusTodo,
		Priority:    issue.PriorityNone,
		Labels:      []label.Label{},
	}
	returned := existing
	returned.Description = &newDesc

	var published []issue.IssueDescriptionUpdatedEvent
	ctx := event.WithPublisher(context.Background(), func(_ context.Context, e issue.IssueDescriptionUpdatedEvent) error {
		published = append(published, e)
		return nil
	})

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)
	repository.On("Update", mock.Anything, mock.MatchedBy(func(i issue.Issue) bool {
		return i.ID == issueID && i.Description == &newDesc
	})).Return(returned, nil)

	actual, err := service.UpdateIssueDescription(ctx, issueID, &newDesc)
	require.NoError(t, err)
	require.Len(t, published, 1)
	assert.Equal(t, returned, actual)
	assert.Equal(t, returned, published[0].Payload)
	repository.AssertExpectations(t)
}

func Test_UpdateIssueDescription_EmitDescriptionUpdatedError_ReturnsError(t *testing.T) {
	repository := &mockIssueRepository{}
	uow := &fakeUnitOfWork{repositories: issue.Repositories{Issues: repository}}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

	issueID := uuid.New()
	newDesc := "new description"
	existing := issue.Issue{
		ID:       issueID,
		Title:    "Test issue",
		Status:   issue.StatusTodo,
		Priority: issue.PriorityNone,
		Labels:   []label.Label{},
	}
	returned := existing
	returned.Description = &newDesc

	expectedError := errors.New("pusblish error")
	ctx := event.WithPublisher(context.Background(), func(_ context.Context, _ issue.IssueDescriptionUpdatedEvent) error {
		return expectedError
	})

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)
	repository.On("Update", mock.Anything, mock.Anything).Return(returned, nil)

	_, err := service.UpdateIssueDescription(ctx, issueID, &newDesc)
	assert.ErrorIs(t, err, expectedError)
	repository.AssertExpectations(t)
}

// — UpdateIssueTitle —

func Test_UpdateIssueTitle_ValidTitle_ReturnsUpdatedIssue(t *testing.T) {
	repository := &mockIssueRepository{}
	uow := &fakeUnitOfWork{repositories: issue.Repositories{Issues: repository}}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

	issueID := uuid.New()
	title := "new title"
	existing := issue.Issue{
		ID:       issueID,
		Title:    "old title",
		Status:   issue.StatusTodo,
		Priority: issue.PriorityNone,
		Labels:   []label.Label{},
	}
	returned := existing
	returned.Title = title

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)
	repository.On("Update", mock.Anything, mock.MatchedBy(func(i issue.Issue) bool {
		return i.ID == issueID && i.Title == title
	})).Return(returned, nil)

	ctxWithNoopPublisher := event.WithPublisher(context.Background(), func(_ context.Context, _ issue.IssueTitleUpdatedEvent) error {
		return nil
	})

	actual, err := service.UpdateIssueTitle(ctxWithNoopPublisher, issueID, title)
	require.NoError(t, err)
	assert.Equal(t, title, actual.Title)
	repository.AssertExpectations(t)
}

func Test_UpdateIssueTitle_EmptyTitle_ReturnsError(t *testing.T) {
	uow := &mockUnitOfWork{}
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

	issueID := uuid.New()
	existing := issue.Issue{
		ID:       issueID,
		Title:    "old title",
		Status:   issue.StatusTodo,
		Priority: issue.PriorityNone,
		Labels:   []label.Label{},
	}

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)

	_, err := service.UpdateIssueTitle(context.Background(), issueID, "")
	require.Error(t, err)
	assert.ErrorIs(t, err, issue.ErrInvalidIssue)
	repository.AssertExpectations(t)
}

func Test_UpdateIssueTitle_IssueNotFound_ReturnsError(t *testing.T) {
	uow := &mockUnitOfWork{}
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

	issueID := uuid.New()
	repository.On("GetIssue", mock.Anything, issueID).Return(issue.Issue{}, issue.ErrIssueNotFound)

	_, err := service.UpdateIssueTitle(context.Background(), issueID, "new title")
	require.Error(t, err)
	assert.ErrorIs(t, err, issue.ErrIssueNotFound)
	repository.AssertExpectations(t)
}

func Test_UpdateIssueTitle_UpdateError_ReturnsError(t *testing.T) {
	repository := &mockIssueRepository{}
	uow := &fakeUnitOfWork{repositories: issue.Repositories{Issues: repository}}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

	issueID := uuid.New()
	existing := issue.Issue{ID: issueID, Title: "old title", Status: issue.StatusTodo, Priority: issue.PriorityNone, Labels: []label.Label{}}
	updateErr := errors.New("db error")

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)
	repository.On("Update", mock.Anything, mock.Anything).Return(issue.Issue{}, updateErr)

	_, err := service.UpdateIssueTitle(context.Background(), issueID, "new title")
	require.Error(t, err)
	assert.ErrorIs(t, err, updateErr)
	repository.AssertExpectations(t)
}

// — UpdateIssuePriority —

func Test_UpdateIssuePriority_ValidPriority_ReturnsUpdatedIssue(t *testing.T) {
	repository := &mockIssueRepository{}
	uow := &fakeUnitOfWork{repositories: issue.Repositories{Issues: repository}}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

	issueID := uuid.New()
	existing := issue.Issue{
		ID:       issueID,
		Title:    "Test issue",
		Status:   issue.StatusTodo,
		Priority: issue.PriorityNone,
		Labels:   []label.Label{},
	}
	returned := existing
	returned.Priority = issue.PriorityHigh

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)
	repository.On("Update", mock.Anything, mock.MatchedBy(func(i issue.Issue) bool {
		return i.ID == issueID && i.Priority == issue.PriorityHigh
	})).Return(returned, nil)

	ctxWithNoopPublisher := event.WithPublisher(context.Background(), func(_ context.Context, _ issue.IssuePriorityUpdatedEvent) error {
		return nil
	})
	actual, err := service.UpdateIssuePriority(ctxWithNoopPublisher, issueID, issue.PriorityHigh)
	require.NoError(t, err)
	assert.Equal(t, issue.PriorityHigh, actual.Priority)
	repository.AssertExpectations(t)
}

func Test_UpdateIssuePriority_InvalidPriority_ReturnsError(t *testing.T) {
	repository := &mockIssueRepository{}
	uow := &fakeUnitOfWork{repositories: issue.Repositories{Issues: repository}}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

	issueID := uuid.New()
	existing := issue.Issue{
		ID:       issueID,
		Title:    "Test issue",
		Status:   issue.StatusTodo,
		Priority: issue.PriorityNone,
		Labels:   []label.Label{},
	}

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)

	_, err := service.UpdateIssuePriority(context.Background(), issueID, issue.Priority("extreme"))
	require.Error(t, err)
	assert.ErrorIs(t, err, issue.ErrInvalidIssue)
	repository.AssertExpectations(t)
}

func Test_UpdateIssuePriority_IssueNotFound_ReturnsError(t *testing.T) {
	repository := &mockIssueRepository{}
	uow := &fakeUnitOfWork{repositories: issue.Repositories{Issues: repository}}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

	issueID := uuid.New()
	repository.On("GetIssue", mock.Anything, issueID).Return(issue.Issue{}, issue.ErrIssueNotFound)

	_, err := service.UpdateIssuePriority(context.Background(), issueID, issue.PriorityHigh)
	require.Error(t, err)
	assert.ErrorIs(t, err, issue.ErrIssueNotFound)
	repository.AssertExpectations(t)
}

func Test_UpdateIssuePriority_UpdateError_ReturnsError(t *testing.T) {
	repository := &mockIssueRepository{}
	uow := &fakeUnitOfWork{repositories: issue.Repositories{Issues: repository}}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

	issueID := uuid.New()
	existing := issue.Issue{ID: issueID, Status: issue.StatusTodo, Priority: issue.PriorityNone, Labels: []label.Label{}}
	updateErr := errors.New("db error")

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)
	repository.On("Update", mock.Anything, mock.Anything).Return(issue.Issue{}, updateErr)

	_, err := service.UpdateIssuePriority(context.Background(), issueID, issue.PriorityHigh)
	require.Error(t, err)
	assert.ErrorIs(t, err, updateErr)
	repository.AssertExpectations(t)
}

func Test_UpdateIssuePriority_SuccessfulUpdate_PublishesIssuePriorityUpdatedEvent(t *testing.T) {
	repository := &mockIssueRepository{}
	uow := &fakeUnitOfWork{repositories: issue.Repositories{Issues: repository}}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

	issueID := uuid.New()
	existing := issue.Issue{
		ID:       issueID,
		Title:    "Test issue",
		Status:   issue.StatusTodo,
		Priority: issue.PriorityNone,
		Labels:   []label.Label{},
	}
	returned := existing
	returned.Priority = issue.PriorityHigh

	var published []issue.IssuePriorityUpdatedEvent
	ctx := event.WithPublisher(context.Background(), func(_ context.Context, e issue.IssuePriorityUpdatedEvent) error {
		published = append(published, e)
		return nil
	})

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)
	repository.On("Update", mock.Anything, mock.MatchedBy(func(i issue.Issue) bool {
		return i.ID == issueID && i.Priority == issue.PriorityHigh
	})).Return(returned, nil)

	actual, err := service.UpdateIssuePriority(ctx, issueID, issue.PriorityHigh)
	require.NoError(t, err)
	require.Len(t, published, 1)
	assert.Equal(t, returned, actual)
	assert.Equal(t, returned, published[0].Payload)
	repository.AssertExpectations(t)
}

func Test_UpdateIssuePriority_EmitPriorityUpdatedError_ReturnsError(t *testing.T) {
	repository := &mockIssueRepository{}
	uow := &fakeUnitOfWork{repositories: issue.Repositories{Issues: repository}}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

	issueID := uuid.New()
	existing := issue.Issue{
		ID:       issueID,
		Title:    "Test issue",
		Status:   issue.StatusTodo,
		Priority: issue.PriorityNone,
		Labels:   []label.Label{},
	}
	returned := existing
	returned.Priority = issue.PriorityHigh

	expectedError := errors.New("nats down")
	ctx := event.WithPublisher(context.Background(), func(_ context.Context, _ issue.IssuePriorityUpdatedEvent) error {
		return expectedError
	})

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)
	repository.On("Update", mock.Anything, mock.Anything).Return(returned, nil)

	_, err := service.UpdateIssuePriority(ctx, issueID, issue.PriorityHigh)
	require.ErrorIs(t, err, expectedError)
	repository.AssertExpectations(t)
}

// — UpdateIssueStatus —

func Test_UpdateIssueStatus_ValidStatus_ReturnsUpdatedIssue(t *testing.T) {
	uow := &mockUnitOfWork{}
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

	issueID := uuid.New()
	existing := issue.Issue{
		ID:       issueID,
		Title:    "Test issue",
		Status:   issue.StatusTodo,
		Priority: issue.PriorityNone,
		Labels:   []label.Label{},
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
	uow := &mockUnitOfWork{}
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

	issueID := uuid.New()
	existing := issue.Issue{
		ID:       issueID,
		Title:    "Test issue",
		Status:   issue.StatusTodo,
		Priority: issue.PriorityNone,
		Labels:   []label.Label{},
	}

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)

	_, err := service.UpdateIssueStatus(context.Background(), issueID, issue.Status("archived"))
	require.Error(t, err)
	assert.ErrorIs(t, err, issue.ErrInvalidIssue)
	repository.AssertExpectations(t)
}

func Test_UpdateIssueStatus_IssueNotFound_ReturnsError(t *testing.T) {
	uow := &mockUnitOfWork{}
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

	issueID := uuid.New()
	repository.On("GetIssue", mock.Anything, issueID).Return(issue.Issue{}, issue.ErrIssueNotFound)

	_, err := service.UpdateIssueStatus(context.Background(), issueID, issue.StatusDone)
	require.Error(t, err)
	assert.ErrorIs(t, err, issue.ErrIssueNotFound)
	repository.AssertExpectations(t)
}

func Test_UpdateIssueStatus_UpdateError_ReturnsError(t *testing.T) {
	uow := &mockUnitOfWork{}
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

	issueID := uuid.New()
	existing := issue.Issue{ID: issueID, Status: issue.StatusTodo, Priority: issue.PriorityNone, Labels: []label.Label{}}
	updateErr := errors.New("db error")

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)
	repository.On("Update", mock.Anything, mock.Anything).Return(issue.Issue{}, updateErr)

	_, err := service.UpdateIssueStatus(context.Background(), issueID, issue.StatusDone)
	require.Error(t, err)
	assert.ErrorIs(t, err, updateErr)
	repository.AssertExpectations(t)
}

func Test_CreateIssue_SuccessfulPersistence_PublishesIssueCreatedEvent(t *testing.T) {
	repository := &mockIssueRepository{}
	labelRepository := &mockLabelRepository{}
	uow := &fakeUnitOfWork{repositories: issue.Repositories{Issues: repository}}
	service := issue.NewIssueService(uow, repository, labelRepository)

	projectID := uuid.New()
	reporterID := uuid.New()
	command := issue.CreateIssueCommand{
		ProjectID:  projectID,
		ReporterID: reporterID,
		Title:      "Event test",
		Status:     issue.StatusTodo,
		Priority:   issue.PriorityMedium,
	}
	expected := issue.Issue{
		ID:         uuid.New(),
		ProjectID:  projectID,
		ReporterID: reporterID,
		Title:      "Event test",
		Status:     issue.StatusTodo,
		Priority:   issue.PriorityMedium,
	}

	var published []issue.IssueCreatedEvent
	ctx := event.WithPublisher(context.Background(), func(_ context.Context, e issue.IssueCreatedEvent) error {
		published = append(published, e)
		return nil
	})

	labelRepository.On("ListByIDs", mock.Anything, command.LabelIDs).Return([]label.Label{}, nil)
	repository.On("CreateIssue", mock.Anything, mock.Anything).Return(expected, nil)
	actual, err := service.CreateIssue(ctx, command)
	require.NoError(t, err)

	require.Len(t, published, 1)
	assert.Equal(t, expected, actual)
	assert.Equal(t, expected, published[0].Payload)
	labelRepository.AssertExpectations(t)
	repository.AssertExpectations(t)
}

func Test_CreateIssue_EmitCreatedError_ReturnsError(t *testing.T) {
	repository := &mockIssueRepository{}
	labelRepository := &mockLabelRepository{}
	uow := &fakeUnitOfWork{repositories: issue.Repositories{Issues: repository}}
	service := issue.NewIssueService(uow, repository, labelRepository)

	command := issue.CreateIssueCommand{
		ProjectID:  uuid.New(),
		ReporterID: uuid.New(),
		Title:      "Publisher fails",
		Status:     issue.StatusBacklog,
		Priority:   issue.PriorityLow,
	}
	returned := issue.Issue{
		ID:    uuid.New(),
		Title: "Publisher fails",
	}

	expectedError := errors.New("publisher down")
	ctx := event.WithPublisher(context.Background(), func(_ context.Context, _ issue.IssueCreatedEvent) error {
		return expectedError
	})

	labelRepository.On("ListByIDs", mock.Anything, command.LabelIDs).Return([]label.Label{}, nil)
	repository.On("CreateIssue", mock.Anything, mock.Anything).Return(returned, nil)

	_, err := service.CreateIssue(ctx, command)
	require.ErrorIs(t, err, expectedError)
	labelRepository.AssertExpectations(t)
	repository.AssertExpectations(t)
}

func Test_UpdateIssueStatus_SuccessfulUpdate_PublishesIssueStatusUpdatedEvent(t *testing.T) {
	uow := &mockUnitOfWork{}
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

	issueID := uuid.New()
	existing := issue.Issue{
		ID:       issueID,
		Title:    "Test issue",
		Status:   issue.StatusTodo,
		Priority: issue.PriorityNone,
		Labels:   []label.Label{},
	}
	returned := existing
	returned.Status = issue.StatusDone

	var published []issue.IssueStatusUpdatedEvent
	ctx := event.WithPublisher(context.Background(), func(_ context.Context, e issue.IssueStatusUpdatedEvent) error {
		published = append(published, e)
		return nil
	})

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)
	repository.On("Update", mock.Anything, mock.MatchedBy(func(i issue.Issue) bool {
		return i.ID == issueID && i.Status == issue.StatusDone
	})).Return(returned, nil)

	actual, err := service.UpdateIssueStatus(ctx, issueID, issue.StatusDone)
	require.NoError(t, err)
	require.Len(t, published, 1)
	assert.Equal(t, returned, actual)
	assert.Equal(t, returned, published[0].Payload)
	repository.AssertExpectations(t)
}

func Test_UpdateIssueStatus_PublisherError_StillReturnsIssue(t *testing.T) {
	uow := &mockUnitOfWork{}
	repository := &mockIssueRepository{}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

	issueID := uuid.New()
	existing := issue.Issue{
		ID:       issueID,
		Title:    "Test issue",
		Status:   issue.StatusTodo,
		Priority: issue.PriorityNone,
		Labels:   []label.Label{},
	}
	returned := existing
	returned.Status = issue.StatusDone

	ctx := event.WithPublisher(context.Background(), func(_ context.Context, _ issue.IssueStatusUpdatedEvent) error {
		return errors.New("nats down")
	})

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)
	repository.On("Update", mock.Anything, mock.Anything).Return(returned, nil)

	actual, err := service.UpdateIssueStatus(ctx, issueID, issue.StatusDone)
	require.NoError(t, err)
	assert.Equal(t, returned, actual)
	repository.AssertExpectations(t)
}

func Test_UpdateIssueTitle_SuccessfulUpdate_PublishesIssueTitleUpdatedEvent(t *testing.T) {
	repository := &mockIssueRepository{}
	uow := &fakeUnitOfWork{repositories: issue.Repositories{Issues: repository}}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

	issueID := uuid.New()
	existing := issue.Issue{
		ID:       issueID,
		Title:    "old title",
		Status:   issue.StatusTodo,
		Priority: issue.PriorityNone,
		Labels:   []label.Label{},
	}
	returned := existing
	returned.Title = "new title"

	var published []issue.IssueTitleUpdatedEvent
	ctx := event.WithPublisher(context.Background(), func(_ context.Context, e issue.IssueTitleUpdatedEvent) error {
		published = append(published, e)
		return nil
	})

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)
	repository.On("Update", mock.Anything, mock.MatchedBy(func(i issue.Issue) bool {
		return i.ID == issueID && i.Title == "new title"
	})).Return(returned, nil)

	actual, err := service.UpdateIssueTitle(ctx, issueID, "new title")
	require.NoError(t, err)
	require.Len(t, published, 1)
	assert.Equal(t, returned, actual)
	assert.Equal(t, returned, published[0].Payload)
	repository.AssertExpectations(t)
}

func Test_UpdateIssueTitle_EmitTitleUpdatedError_ReturnsError(t *testing.T) {
	repository := &mockIssueRepository{}
	uow := &fakeUnitOfWork{repositories: issue.Repositories{Issues: repository}}
	service := issue.NewIssueService(uow, repository, &mockLabelRepository{})

	issueID := uuid.New()
	existing := issue.Issue{
		ID:       issueID,
		Title:    "old title",
		Status:   issue.StatusTodo,
		Priority: issue.PriorityNone,
		Labels:   []label.Label{},
	}
	returned := existing
	returned.Title = "new title"

	expectedError := errors.New("nats down")
	ctx := event.WithPublisher(context.Background(), func(_ context.Context, _ issue.IssueTitleUpdatedEvent) error {
		return expectedError
	})

	repository.On("GetIssue", mock.Anything, issueID).Return(existing, nil)
	repository.On("Update", mock.Anything, mock.Anything).Return(returned, nil)

	_, err := service.UpdateIssueTitle(ctx, issueID, "new title")
	require.ErrorIs(t, err, expectedError)
	repository.AssertExpectations(t)
}
