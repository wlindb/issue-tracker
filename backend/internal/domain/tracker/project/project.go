package project

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID
	Name        string
	Description *string
	OwnerID     uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ListProjectQuery holds all parameters for listing projects.
type ListProjectQuery struct {
	Cursor *string
	Limit  *int
}

// Projects is the paginated result of a List operation.
type Projects struct {
	Items []Project
}

const defaultLimit = 20

func NewListProjectQuery(cursor *string, limit *int) ListProjectQuery {
	if limit == nil {
		l := defaultLimit
		limit = &l
	}
	return ListProjectQuery{Cursor: cursor, Limit: limit}
}

type ProjectRepository interface {
	Create(ctx context.Context, id uuid.UUID, ownerID uuid.UUID, name string, description *string) (*Project, error)
	List(ctx context.Context, query ListProjectQuery) (Projects, error)
}

var (
	ErrProjectNotFound = errors.New("project not found")
	ErrInvalidProject  = errors.New("invalid project")
)
