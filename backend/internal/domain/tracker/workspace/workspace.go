package workspace

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrWorkspaceNotFound = errors.New("workspace not found")
	ErrInvalidWorkspace  = errors.New("invalid workspace")
)

type Workspace struct {
	ID        uuid.UUID
	Name      string
	OwnerID   uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type WorkspaceMember struct {
	WorkspaceID uuid.UUID
	UserID      uuid.UUID
}

type WorkspaceRepository interface {
	Create(ctx context.Context, id uuid.UUID, ownerID uuid.UUID, name string) (*Workspace, error)
	Get(ctx context.Context, id uuid.UUID) (*Workspace, error)
	List(ctx context.Context, userID uuid.UUID) ([]Workspace, error)
	IsMember(ctx context.Context, workspaceID uuid.UUID, userID uuid.UUID) (bool, error)
}
