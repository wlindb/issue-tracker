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

func (m *mockWorkspaceRepository) Create(ctx context.Context, w workspace.Workspace) (workspace.Workspace, error) {
	args := m.Called(ctx, w)
	return args.Get(0).(workspace.Workspace), args.Error(1)
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

func (m *mockWorkspaceRepository) IsMember(ctx context.Context, workspaceID uuid.UUID, userID uuid.UUID) (bool, error) {
	args := m.Called(ctx, workspaceID, userID)
	return args.Bool(0), args.Error(1)
}

func Test_Create_ValidWorkspace_ReturnsWorkspace(t *testing.T) {
	repository := &mockWorkspaceRepository{}
	ownerID := uuid.New()
	workspaceID := uuid.New()
	now := time.Now().UTC()
	input := workspace.Workspace{
		ID:      workspaceID,
		Name:    "Acme",
		OwnerID: ownerID,
	}
	expected := workspace.Workspace{
		ID:        workspaceID,
		Name:      "Acme",
		OwnerID:   ownerID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	repository.On("Create", mock.Anything, input).Return(expected, nil)

	service := workspace.NewWorkspaceService(repository)
	actual, err := service.Create(context.Background(), input)

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
	repository.AssertExpectations(t)
}

func Test_Create_RepositoryError_ReturnsWrappedError(t *testing.T) {
	repository := &mockWorkspaceRepository{}
	ownerID := uuid.New()
	repoErr := errors.New("db down")
	input := workspace.Workspace{
		ID:      uuid.New(),
		Name:    "Acme",
		OwnerID: ownerID,
	}

	repository.On("Create", mock.Anything, input).Return(workspace.Workspace{}, repoErr)

	service := workspace.NewWorkspaceService(repository)
	_, err := service.Create(context.Background(), input)

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

func Test_IsMember_MemberExists_ReturnsTrue(t *testing.T) {
	repository := &mockWorkspaceRepository{}
	workspaceID, userID := uuid.New(), uuid.New()

	repository.On("IsMember", mock.Anything, workspaceID, userID).Return(true, nil)

	service := workspace.NewWorkspaceService(repository)
	actual, err := service.IsMember(context.Background(), workspaceID, userID)

	require.NoError(t, err)
	assert.True(t, actual)
	repository.AssertExpectations(t)
}

func Test_IsMember_NotAMember_ReturnsFalse(t *testing.T) {
	repository := &mockWorkspaceRepository{}
	workspaceID, userID := uuid.New(), uuid.New()

	repository.On("IsMember", mock.Anything, workspaceID, userID).Return(false, nil)

	service := workspace.NewWorkspaceService(repository)
	actual, err := service.IsMember(context.Background(), workspaceID, userID)

	require.NoError(t, err)
	assert.False(t, actual)
	repository.AssertExpectations(t)
}

func Test_IsMember_RepositoryError_ReturnsWrappedError(t *testing.T) {
	repository := &mockWorkspaceRepository{}
	workspaceID, userID := uuid.New(), uuid.New()
	repoErr := errors.New("db down")

	repository.On("IsMember", mock.Anything, workspaceID, userID).Return(false, repoErr)

	service := workspace.NewWorkspaceService(repository)
	_, err := service.IsMember(context.Background(), workspaceID, userID)

	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
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
