package tracker

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	issuedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
)

var _ issuedomain.LabelRepository = (*TracingLabelRepository)(nil)

// TracingLabelRepository wraps a LabelRepository and adds an OTel child span to each operation.
type TracingLabelRepository struct {
	inner  issuedomain.LabelRepository
	tracer trace.Tracer
}

// NewTracingLabelRepository returns a TracingLabelRepository that delegates to inner.
func NewTracingLabelRepository(inner issuedomain.LabelRepository, tracer trace.Tracer) *TracingLabelRepository {
	return &TracingLabelRepository{inner: inner, tracer: tracer}
}

func (r *TracingLabelRepository) Upsert(ctx context.Context, label issuedomain.Label) (issuedomain.Label, error) {
	ctx, span := r.tracer.Start(ctx, "tracker.LabelRepository.Upsert")
	defer span.End()

	result, err := r.inner.Upsert(ctx, label)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return issuedomain.Label{}, fmt.Errorf("upsert label: %w", err)
	}
	return result, nil
}

func (r *TracingLabelRepository) ListByIDs(ctx context.Context, ids []uuid.UUID) ([]issuedomain.Label, error) {
	ctx, span := r.tracer.Start(ctx, "tracker.LabelRepository.ListByIDs")
	defer span.End()

	result, err := r.inner.ListByIDs(ctx, ids)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return []issuedomain.Label{}, fmt.Errorf("list labels by ids: %w", err)
	}
	return result, nil
}

func (r *TracingLabelRepository) SearchByName(ctx context.Context, name string) ([]issuedomain.Label, error) {
	ctx, span := r.tracer.Start(ctx, "tracker.LabelRepository.SearchByName")
	defer span.End()

	result, err := r.inner.SearchByName(ctx, name)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return []issuedomain.Label{}, fmt.Errorf("search labels by name: %w", err)
	}
	return result, nil
}
