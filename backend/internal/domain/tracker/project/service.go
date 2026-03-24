package project

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type ProjectService struct {
	repository ProjectRepository
}

func NewProjectService(repository ProjectRepository) *ProjectService {
	return &ProjectService{repository: repository}
}

func (s *ProjectService) Create(ctx context.Context, id uuid.UUID, ownerID uuid.UUID, name string, description *string) (*Project, error) {
	if name == "" {
		return nil, fmt.Errorf("%w: name is required", ErrInvalidProject)
	}
	p, err := s.repository.Create(ctx, id, ownerID, name, description)
	if err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}
	return p, nil
}

func (s *ProjectService) List(ctx context.Context, query ListProjectQuery) (Projects, error) {
	projects, err := s.repository.List(ctx, query)
	if err != nil {
		return Projects{}, fmt.Errorf("list projects: %w", err)
	}
	return projects, nil
}
