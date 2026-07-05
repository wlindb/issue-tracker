package label

import (
	"context"

	"github.com/google/uuid"
)

// Label is a workspace-scoped tag that can be attached to issues.
type Label struct {
	ID   uuid.UUID
	Name string
}

// LabelRepository defines the persistence interface for labels.
type LabelRepository interface {
	GetOrCreate(ctx context.Context, name string) (Label, error)
	ListByIDs(ctx context.Context, ids []uuid.UUID) ([]Label, error)
	SearchByName(ctx context.Context, name string) ([]Label, error)
}
