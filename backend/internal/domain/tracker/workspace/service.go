package workspace

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type WorkspaceService struct {
	repository WorkspaceRepository
}

func NewWorkspaceService(repository WorkspaceRepository) *WorkspaceService {
	return &WorkspaceService{repository: repository}
}

func (s *WorkspaceService) Create(ctx context.Context, workspace Workspace) (Workspace, error) {
	w, err := s.repository.Create(ctx, workspace)
	if err != nil {
		return Workspace{}, fmt.Errorf("create workspace: %w", err)
	}
	return w, nil
}

func (s *WorkspaceService) Get(ctx context.Context, id uuid.UUID) (*Workspace, error) {
	w, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get workspace: %w", err)
	}
	return w, nil
}

func (s *WorkspaceService) List(ctx context.Context, userID uuid.UUID) ([]Workspace, error) {
	workspaces, err := s.repository.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list workspaces: %w", err)
	}
	return workspaces, nil
}

func (s *WorkspaceService) IsMember(ctx context.Context, workspaceID uuid.UUID, userID uuid.UUID) (bool, error) {
	member, err := s.repository.IsMember(ctx, workspaceID, userID)
	if err != nil {
		return false, fmt.Errorf("check workspace membership: %w", err)
	}
	return member, nil
}
