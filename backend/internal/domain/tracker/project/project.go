package project

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	nonSlugCharPattern    = regexp.MustCompile(`[^a-z0-9-]`)
	multipleDashesPattern = regexp.MustCompile(`-{2,}`)
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

// CreateProjectCommand holds all inputs needed to create a new project.
type CreateProjectCommand struct {
	Name        string
	Description *string
	OwnerID     uuid.UUID
}

// Slugify converts s into a URL-friendly slug.
func (c CreateProjectCommand) Slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = nonSlugCharPattern.ReplaceAllString(s, "")
	s = multipleDashesPattern.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

// Slugifier is a function that converts a string to a slug.
type Slugifier func(s string) string

// ToProject builds a Project from the command using the given id and slugifier.
func (c CreateProjectCommand) ToProject(id uuid.UUID, slugifier Slugifier) Project {
	now := time.Now()
	return Project{
		ID:          id,
		Identifier:  slugifier(c.Name),
		Name:        c.Name,
		Description: c.Description,
		OwnerID:     c.OwnerID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
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
