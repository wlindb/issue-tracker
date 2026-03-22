package project

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type ProjectService struct {
	repo ProjectRepository
}

func NewProjectService(repo ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) Create(ctx context.Context, id uuid.UUID, ownerID uuid.UUID, name string, description *string) (*Project, error) {
	if name == "" {
		return nil, fmt.Errorf("%w: name is required", ErrInvalidProject)
	}
	p, err := s.repo.Create(ctx, id, ownerID, name, description)
	if err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}
	return p, nil
}

func (s *ProjectService) List(_ context.Context, _ uuid.UUID, _ *string, _ *int) ([]Project, *string, error) {
	return nil, nil, fmt.Errorf("not implemented")
}
