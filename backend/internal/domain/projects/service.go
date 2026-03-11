package projects

import (
	"context"

	"github.com/google/uuid"
)

type ProjectService struct {
	repo ProjectRepository
}

func NewProjectService(repo ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) Create(ctx context.Context, ownerID uuid.UUID, name string, description *string) (*Project, error) {
	return s.repo.Create(ctx, ownerID, name, description)
}
