package tracker

import (
	"context"
	"fmt"
	"math"

	"github.com/jackc/pgx/v5/pgxpool"

	projectdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/project"
	trackerdb "github.com/wlindb/issue-tracker/internal/infrastructure/tracker/generated"
)

// Compile-time: *ProjectRepository must satisfy domain interface.
var _ projectdomain.ProjectRepository = (*ProjectRepository)(nil)

// ProjectRepository is a PostgreSQL-backed implementation of projectdomain.ProjectRepository.
type ProjectRepository struct {
	queries *trackerdb.Queries
}

// NewProjectRepository returns a ProjectRepository backed by the given pool.
func NewProjectRepository(pool *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{queries: trackerdb.New(pool)}
}

// Create inserts a new project row and returns the domain model.
func (r *ProjectRepository) Create(ctx context.Context, project projectdomain.Project) (projectdomain.Project, error) {
	row, err := r.queries.CreateProject(ctx, createProjectParamsFromDomain(project))
	if err != nil {
		return projectdomain.Project{}, fmt.Errorf("create project: %w", err)
	}
	return projectToDomain(row), nil
}

// List returns up to query.Limit projects ordered by created_at descending.
func (r *ProjectRepository) List(ctx context.Context, query projectdomain.ListProjectQuery) (projectdomain.Projects, error) {
	limit := *query.Limit
	if limit < 0 || limit > math.MaxInt32 {
		return projectdomain.Projects{}, fmt.Errorf("list projects: limit out of range: %d", limit)
	}
	rows, err := r.queries.ListProjects(ctx, int32(limit))
	if err != nil {
		return projectdomain.Projects{}, fmt.Errorf("list projects: %w", err)
	}
	return projectdomain.Projects{Items: projectsToDomain(rows)}, nil
}
