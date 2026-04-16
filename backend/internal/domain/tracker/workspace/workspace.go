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

type WorkspaceMembers struct {
	Members []WorkspaceMember
}

// New constructs and validates a Workspace value.
func New(id uuid.UUID, name string, ownerID uuid.UUID) (Workspace, error) {
	if id == uuid.Nil {
		return Workspace{}, ErrInvalidWorkspace
	}
	if name == "" {
		return Workspace{}, ErrInvalidWorkspace
	}
	if ownerID == uuid.Nil {
		return Workspace{}, ErrInvalidWorkspace
	}
	now := time.Now()
	return Workspace{
		ID:        id,
		Name:      name,
		OwnerID:   ownerID,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

type WorkspaceRepository interface {
	Create(ctx context.Context, workspace Workspace) (Workspace, error)
	Get(ctx context.Context, id uuid.UUID) (*Workspace, error)
	List(ctx context.Context, userID uuid.UUID) ([]Workspace, error)
	ListMembers(ctx context.Context, workspaceID uuid.UUID) (WorkspaceMembers, error)
	IsMember(ctx context.Context, workspaceID uuid.UUID, userID uuid.UUID) (bool, error)
}
