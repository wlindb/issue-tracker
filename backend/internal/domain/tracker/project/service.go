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

func (s *ProjectService) Create(ctx context.Context, command CreateProjectCommand) (Project, error) {
	project, err := command.ToProject(uuid.New(), command.Slugify)
	if err != nil {
		return Project{}, fmt.Errorf("create project: %w", err)
	}
	result, err := s.repository.Create(ctx, project)
	if err != nil {
		return Project{}, fmt.Errorf("create project: %w", err)
	}
	return result, nil
}

func (s *ProjectService) List(ctx context.Context, query ListProjectQuery) (Projects, error) {
	projects, err := s.repository.List(ctx, query)
	if err != nil {
		return Projects{}, fmt.Errorf("list projects: %w", err)
	}
	return projects, nil
}
