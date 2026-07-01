package issue

import (
	"context"

	"github.com/google/uuid"
)

type LabelRepository interface {
	Upsert(ctx context.Context, label Label) (Label, error)
	ListByIDs(ctx context.Context, ids []uuid.UUID) ([]Label, error)
	SearchByName(ctx context.Context, name string) ([]Label, error)
}
