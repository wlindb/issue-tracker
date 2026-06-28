package tracker

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	issuedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
	trackerdb "github.com/wlindb/issue-tracker/internal/infrastructure/tracker/generated"
)

var _ issuedomain.LabelRepository = (*LabelRepository)(nil)

// LabelRepository is a PostgreSQL-backed implementation of issuedomain.LabelRepository.
type LabelRepository struct {
	db trackerdb.DBTX
}

// NewLabelRepository returns a LabelRepository backed by the given database connection.
func NewLabelRepository(db trackerdb.DBTX) *LabelRepository {
	return &LabelRepository{db: db}
}

// Upsert inserts a new label or updates its name on ID conflict, returning the domain model.
func (r *LabelRepository) Upsert(ctx context.Context, label issuedomain.Label) (issuedomain.Label, error) {
	queries := trackerdb.New(r.db)
	row, err := queries.UpsertLabel(ctx, upsertLabelParamsFromDomain(label))
	if err != nil {
		return issuedomain.Label{}, fmt.Errorf("upsert label: %w", err)
	}
	return labelToDomain(row), nil
}

// ListByIDs retrieves all labels with the given IDs, scoped to the current workspace via RLS.
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

// SearchByName returns all labels whose name contains the given substring (case-insensitive),
// scoped to the current workspace via RLS.
func (r *LabelRepository) SearchByName(ctx context.Context, name string) ([]issuedomain.Label, error) {
	queries := trackerdb.New(r.db)
	rows, err := queries.SearchLabelsByName(ctx, name)
	if err != nil {
		return []issuedomain.Label{}, fmt.Errorf("search labels by name: %w", err)
	}
	if len(rows) == 0 {
		return []issuedomain.Label{}, nil
	}
	return labelsToDomain(rows), nil
}
