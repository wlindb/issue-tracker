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

func TestProjectService_Create_Success(t *testing.T) {
	repo := &mockProjectRepository{}
	svc := project.NewProjectService(repo)

	id := uuid.New()
	ownerID := uuid.New()
	expected := &project.Project{ID: id, OwnerID: ownerID, Name: "My Project"}

	repo.On("Create", mock.Anything, id, ownerID, "My Project", (*string)(nil)).
		Return(expected, nil)

	got, err := svc.Create(context.Background(), id, ownerID, "My Project", nil)
	require.NoError(t, err)
	assert.Equal(t, expected, got)
	repo.AssertExpectations(t)
}

func TestProjectService_Create_WithDescription(t *testing.T) {
	repo := &mockProjectRepository{}
	svc := project.NewProjectService(repo)

	id := uuid.New()
	ownerID := uuid.New()
	desc := "A description"
	expected := &project.Project{ID: id, OwnerID: ownerID, Name: "My Project", Description: &desc}

	repo.On("Create", mock.Anything, id, ownerID, "My Project", &desc).
		Return(expected, nil)

	got, err := svc.Create(context.Background(), id, ownerID, "My Project", &desc)
	require.NoError(t, err)
	assert.Equal(t, expected, got)
	repo.AssertExpectations(t)
}

func TestProjectService_Create_EmptyName(t *testing.T) {
	repo := &mockProjectRepository{}
	svc := project.NewProjectService(repo)

	_, err := svc.Create(context.Background(), uuid.New(), uuid.New(), "", nil)
	require.Error(t, err)
	assert.ErrorIs(t, err, project.ErrInvalidProject)
	repo.AssertNotCalled(t, "Create")
}

func TestProjectService_Create_RepoError(t *testing.T) {
	repo := &mockProjectRepository{}
	svc := project.NewProjectService(repo)

	id := uuid.New()
	ownerID := uuid.New()
	repoErr := errors.New("db error")

	repo.On("Create", mock.Anything, id, ownerID, "My Project", (*string)(nil)).
		Return(nil, repoErr)

	_, err := svc.Create(context.Background(), id, ownerID, "My Project", nil)
	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	repo.AssertExpectations(t)
}
