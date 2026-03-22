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
	OwnerID uuid.UUID
	Cursor  *string
	Limit   *int
}

// Projects is the paginated result of a List operation.
type Projects struct {
	Items      []Project
	nextCursor *string
}

// NewProjects constructs a Projects result with the given items and next-page cursor.
func NewProjects(items []Project, cursor *string) Projects {
	return Projects{Items: items, nextCursor: cursor}
}

// Cursor returns the next page cursor, or nil when there are no more pages.
func (p Projects) Cursor() *string {
	return p.nextCursor
}

type ProjectRepository interface {
	Create(ctx context.Context, id uuid.UUID, ownerID uuid.UUID, name string, description *string) (*Project, error)
}

var ErrProjectNotFound = errors.New("project not found")
var ErrInvalidProject = errors.New("invalid project")
