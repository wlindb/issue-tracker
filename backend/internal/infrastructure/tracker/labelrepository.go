package tracker

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	issuedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
	trackerdb "github.com/wlindb/issue-tracker/internal/infrastructure/tracker/generated"
)

var _ issuedomain.LabelRepository = (*LabelRepository)(nil)

type LabelRepository struct {
	db trackerdb.DBTX
}

func NewLabelRepository(db trackerdb.DBTX) *LabelRepository {
	return &LabelRepository{db: db}
}

func (r *LabelRepository) GetOrCreate(ctx context.Context, name string) (issuedomain.Label, error) {
	queries := trackerdb.New(r.db)
	row, err := queries.GetOrCreateLabel(ctx, getOrCreateLabelParams(uuid.New(), name))
	if err != nil {
		return issuedomain.Label{}, fmt.Errorf("get or create label: %w", err)
	}
	return labelToDomain(row), nil
}

func (r *LabelRepository) ListByIDs(ctx context.Context, ids []uuid.UUID) ([]issuedomain.Label, error) {
	queries := trackerdb.New(r.db)
	rows, err := queries.ListLabelsByIDs(ctx, ids)
	if err != nil {
		return []issuedomain.Label{}, fmt.Errorf("list labels by ids: %w", err)
	}
	if len(rows) == 0 {
		return []issuedomain.Label{}, nil
	}
	return labelsToDomain(rows), nil
}

func (r *LabelRepository) SearchByName(ctx context.Context, name string) ([]issuedomain.Label, error) {
	queries := trackerdb.New(r.db)
	rows, err := queries.SearchLabelsByName(ctx, pgtype.Text{String: name, Valid: true})
	if err != nil {
		return []issuedomain.Label{}, fmt.Errorf("search labels by name: %w", err)
	}
	if len(rows) == 0 {
		return []issuedomain.Label{}, nil
	}
	return labelsToDomain(rows), nil
}
