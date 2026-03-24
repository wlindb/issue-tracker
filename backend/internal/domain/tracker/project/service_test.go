package project_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/wlindb/issue-tracker/internal/domain/tracker/project"
)

type mockProjectRepository struct {
	mock.Mock
}

func (m *mockProjectRepository) Create(ctx context.Context, id uuid.UUID, ownerID uuid.UUID, name string, description *string) (*project.Project, error) {
	args := m.Called(ctx, id, ownerID, name, description)
	if p, ok := args.Get(0).(*project.Project); ok {
		return p, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockProjectRepository) List(ctx context.Context, query project.ListProjectQuery) (project.Projects, error) {
	args := m.Called(ctx, query)
	projects, _ := args.Get(0).(project.Projects)
	return projects, args.Error(1)
}

func TestProjectService_Create_Success(t *testing.T) {
	repository := &mockProjectRepository{}
	service := project.NewProjectService(repository)

	id := uuid.New()
	ownerID := uuid.New()
	expected := &project.Project{ID: id, OwnerID: ownerID, Name: "My Project"}

	repository.On("Create", mock.Anything, id, ownerID, "My Project", (*string)(nil)).
		Return(expected, nil)

	got, err := service.Create(context.Background(), id, ownerID, "My Project", nil)
	require.NoError(t, err)
	assert.Equal(t, expected, got)
	repository.AssertExpectations(t)
}

func TestProjectService_Create_WithDescription(t *testing.T) {
	repository := &mockProjectRepository{}
	service := project.NewProjectService(repository)

	id := uuid.New()
	ownerID := uuid.New()
	description := "A description"
	expected := &project.Project{ID: id, OwnerID: ownerID, Name: "My Project", Description: &description}

	repository.On("Create", mock.Anything, id, ownerID, "My Project", &description).
		Return(expected, nil)

	got, err := service.Create(context.Background(), id, ownerID, "My Project", &description)
	require.NoError(t, err)
	assert.Equal(t, expected, got)
	repository.AssertExpectations(t)
}

func TestProjectService_Create_EmptyName(t *testing.T) {
	repository := &mockProjectRepository{}
	service := project.NewProjectService(repository)

	_, err := service.Create(context.Background(), uuid.New(), uuid.New(), "", nil)
	require.Error(t, err)
	assert.ErrorIs(t, err, project.ErrInvalidProject)
	repository.AssertNotCalled(t, "Create")
}

func TestProjectService_Create_RepoError(t *testing.T) {
	repository := &mockProjectRepository{}
	service := project.NewProjectService(repository)

	id := uuid.New()
	ownerID := uuid.New()
	repoErr := errors.New("db error")

	repository.On("Create", mock.Anything, id, ownerID, "My Project", (*string)(nil)).
		Return(nil, repoErr)

	_, err := service.Create(context.Background(), id, ownerID, "My Project", nil)
	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	repository.AssertExpectations(t)
}

func Test_List_WithProjects_ReturnsProjects(t *testing.T) {
	repository := &mockProjectRepository{}
	service := project.NewProjectService(repository)

	query := project.NewListProjectQuery(nil, nil)
	expected := project.Projects{
		Items: []project.Project{
			{ID: uuid.New(), Name: "Alpha"},
			{ID: uuid.New(), Name: "Beta"},
		},
	}

	repository.On("List", mock.Anything, query).Return(expected, nil)

	actual, err := service.List(context.Background(), query)
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
	repository.AssertExpectations(t)
}

func Test_List_RepoError_ReturnsError(t *testing.T) {
	repository := &mockProjectRepository{}
	service := project.NewProjectService(repository)

	query := project.NewListProjectQuery(nil, nil)
	repoErr := errors.New("db error")

	repository.On("List", mock.Anything, query).Return(project.Projects{}, repoErr)

	_, err := service.List(context.Background(), query)
	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	repository.AssertExpectations(t)
}
