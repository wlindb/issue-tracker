//go:build !integration

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

func (m *mockProjectRepository) Create(ctx context.Context, p project.Project) (project.Project, error) {
	args := m.Called(ctx, p)
	return args.Get(0).(project.Project), args.Error(1)
}

func (m *mockProjectRepository) List(ctx context.Context, query project.ListProjectQuery) (project.Projects, error) {
	args := m.Called(ctx, query)
	projects, _ := args.Get(0).(project.Projects)
	return projects, args.Error(1)
}

func Test_Create_ValidProject_ReturnsProject(t *testing.T) {
	repository := &mockProjectRepository{}
	service := project.NewProjectService(repository)

	ownerID := uuid.New()
	command := project.CreateProjectCommand{Name: "My Project", OwnerID: ownerID}
	expected := project.Project{ID: uuid.New(), Identifier: "my-project", OwnerID: ownerID, Name: "My Project"}

	repository.On("Create", mock.Anything, mock.MatchedBy(func(p project.Project) bool {
		return p.Name == command.Name && p.OwnerID == command.OwnerID && p.Identifier != ""
	})).Return(expected, nil)

	actual, err := service.Create(context.Background(), command)
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
	repository.AssertExpectations(t)
}

func Test_Create_WithDescription_ReturnsProject(t *testing.T) {
	repository := &mockProjectRepository{}
	service := project.NewProjectService(repository)

	ownerID := uuid.New()
	description := "A description"
	command := project.CreateProjectCommand{Name: "My Project", Description: &description, OwnerID: ownerID}
	expected := project.Project{ID: uuid.New(), Identifier: "my-project", OwnerID: ownerID, Name: "My Project", Description: &description}

	repository.On("Create", mock.Anything, mock.MatchedBy(func(p project.Project) bool {
		return p.Name == command.Name && p.OwnerID == command.OwnerID
	})).Return(expected, nil)

	actual, err := service.Create(context.Background(), command)
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
	repository.AssertExpectations(t)
}

func Test_Create_RepositoryError_ReturnsError(t *testing.T) {
	repository := &mockProjectRepository{}
	service := project.NewProjectService(repository)

	command := project.CreateProjectCommand{Name: "My Project", OwnerID: uuid.New()}
	repoErr := errors.New("db error")

	repository.On("Create", mock.Anything, mock.MatchedBy(func(p project.Project) bool {
		return p.Name == command.Name
	})).Return(project.Project{}, repoErr)

	_, err := service.Create(context.Background(), command)
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
