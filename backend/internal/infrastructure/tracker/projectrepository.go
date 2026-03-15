package tracker

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	projectdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/project"
	trackerdb "github.com/wlindb/issue-tracker/internal/infrastructure/tracker/generated"
)

// ProjectRepository is a PostgreSQL-backed implementation of projectdomain.ProjectRepository.
type ProjectRepository struct {
	q *trackerdb.Queries
}

// NewProjectRepository returns a ProjectRepository backed by the given pool.
func NewProjectRepository(pool *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{q: trackerdb.New(pool)}
}

// Create inserts a new project row and returns the domain model.
func (r *ProjectRepository) Create(ctx context.Context, id, ownerID uuid.UUID, name string, description *string) (*projectdomain.Project, error) {
	var desc pgtype.Text
	if description != nil {
		desc = pgtype.Text{String: *description, Valid: true}
	}
	row, err := r.q.CreateProject(ctx, trackerdb.CreateProjectParams{
		ID:          id,
		OwnerID:     ownerID,
		Name:        name,
		Description: desc,
	})
	if err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}
	return rowToProject(row), nil
}

func rowToProject(row trackerdb.Project) *projectdomain.Project {
	p := &projectdomain.Project{
		ID:        row.ID,
		OwnerID:   row.OwnerID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
	if row.Description.Valid {
		s := row.Description.String
		p.Description = &s
	}
	return p
}
