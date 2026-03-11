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

type ProjectRepository interface {
	Create(ctx context.Context, ownerID uuid.UUID, name string, description *string) (*Project, error)
}

var ErrProjectNotFound = errors.New("project not found")
