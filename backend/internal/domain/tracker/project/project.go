package project

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID
	Identifier  string
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

func New(id uuid.UUID, identifier string, name string, description *string, ownerID uuid.UUID) (Project, error) {
	if id == uuid.Nil {
		return Project{}, ErrInvalidProject
	}
	if identifier == "" {
		return Project{}, ErrInvalidProject
	}
	if name == "" {
		return Project{}, ErrInvalidProject
	}
	if ownerID == uuid.Nil {
		return Project{}, ErrInvalidProject
	}
	now := time.Now()
	return Project{
		ID:          id,
		Identifier:  identifier,
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

type ProjectRepository interface {
	Create(ctx context.Context, project Project) (Project, error)
	List(ctx context.Context, query ListProjectQuery) (Projects, error)
}

var (
	ErrProjectNotFound = errors.New("project not found")
	ErrInvalidProject  = errors.New("invalid project")
)
