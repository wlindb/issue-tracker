//go:build !integration

package workspace_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/wlindb/issue-tracker/internal/domain/tracker/workspace"
)

type mockWorkspaceRepository struct {
	mock.Mock
}

func (m *mockWorkspaceRepository) Create(ctx context.Context, id uuid.UUID, ownerID uuid.UUID, name string) (*workspace.Workspace, error) {
	args := m.Called(ctx, id, ownerID, name)
	if w, ok := args.Get(0).(*workspace.Workspace); ok {
		return w, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockWorkspaceRepository) Get(ctx context.Context, id uuid.UUID) (*workspace.Workspace, error) {
	args := m.Called(ctx, id)
	if w, ok := args.Get(0).(*workspace.Workspace); ok {
		return w, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockWorkspaceRepository) List(ctx context.Context, userID uuid.UUID) ([]workspace.Workspace, error) {
	args := m.Called(ctx, userID)
	workspaces, _ := args.Get(0).([]workspace.Workspace)
	return workspaces, args.Error(1)
}

func Test_Create_ValidName_ReturnsWorkspace(t *testing.T) {
	repository := &mockWorkspaceRepository{}
	ownerID := uuid.New()
	now := time.Now().UTC()
	expected := &workspace.Workspace{
		ID:        uuid.New(),
		Name:      "Acme",
		OwnerID:   ownerID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	repository.On("Create", mock.Anything, mock.AnythingOfType("uuid.UUID"), ownerID, "Acme").
		Return(expected, nil)

	service := workspace.NewWorkspaceService(repository)
	actual, err := service.Create(context.Background(), ownerID, "Acme")

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
	repository.AssertExpectations(t)
}

func Test_Create_EmptyName_ReturnsErrInvalidWorkspace(t *testing.T) {
	repository := &mockWorkspaceRepository{}
	service := workspace.NewWorkspaceService(repository)

	_, err := service.Create(context.Background(), uuid.New(), "")

	require.Error(t, err)
	assert.ErrorIs(t, err, workspace.ErrInvalidWorkspace)
	repository.AssertNotCalled(t, "Create")
}

func Test_Create_RepositoryError_ReturnsWrappedError(t *testing.T) {
	repository := &mockWorkspaceRepository{}
	ownerID := uuid.New()
	repoErr := errors.New("db down")

	repository.On("Create", mock.Anything, mock.AnythingOfType("uuid.UUID"), ownerID, "Acme").
		Return(nil, repoErr)

	service := workspace.NewWorkspaceService(repository)
	_, err := service.Create(context.Background(), ownerID, "Acme")

	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	repository.AssertExpectations(t)
}

func Test_Get_ExistingWorkspace_ReturnsWorkspace(t *testing.T) {
	repository := &mockWorkspaceRepository{}
	workspaceID := uuid.New()
	now := time.Now().UTC()
	expected := &workspace.Workspace{
		ID:        workspaceID,
		Name:      "Acme",
		OwnerID:   uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
	}

	repository.On("Get", mock.Anything, workspaceID).Return(expected, nil)

	service := workspace.NewWorkspaceService(repository)
	actual, err := service.Get(context.Background(), workspaceID)

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
	repository.AssertExpectations(t)
}

func Test_Get_RepositoryReturnsNotFound_ReturnsWrappedError(t *testing.T) {
	repository := &mockWorkspaceRepository{}
	workspaceID := uuid.New()

	repository.On("Get", mock.Anything, workspaceID).
		Return(nil, workspace.ErrWorkspaceNotFound)

	service := workspace.NewWorkspaceService(repository)
	_, err := service.Get(context.Background(), workspaceID)

	require.Error(t, err)
	assert.ErrorIs(t, err, workspace.ErrWorkspaceNotFound)
	repository.AssertExpectations(t)
}

func Test_List_ExistingUser_ReturnsWorkspaces(t *testing.T) {
	repository := &mockWorkspaceRepository{}
	userID := uuid.New()
	now := time.Now().UTC()
	expected := []workspace.Workspace{
		{ID: uuid.New(), Name: "Acme", OwnerID: userID, CreatedAt: now, UpdatedAt: now},
	}

	repository.On("List", mock.Anything, userID).Return(expected, nil)

	service := workspace.NewWorkspaceService(repository)
	actual, err := service.List(context.Background(), userID)

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
	repository.AssertExpectations(t)
}

func Test_List_RepositoryError_ReturnsWrappedError(t *testing.T) {
	repository := &mockWorkspaceRepository{}
	userID := uuid.New()
	repoErr := errors.New("db down")

	repository.On("List", mock.Anything, userID).Return(nil, repoErr)

	service := workspace.NewWorkspaceService(repository)
	_, err := service.List(context.Background(), userID)

	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	repository.AssertExpectations(t)
}
