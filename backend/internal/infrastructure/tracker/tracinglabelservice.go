package tracker

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/wlindb/issue-tracker/internal/domain/tracker/label"
)

// labelServicer mirrors api.LabelServicer without importing the api package,
// avoiding a layering violation (infrastructure must not depend on api).
type labelServicer interface {
	Create(ctx context.Context, name string) (label.Label, error)
	Search(ctx context.Context, name string) ([]label.Label, error)
}

// TracingLabelService wraps a labelServicer and adds an OTel child span to each operation.
type TracingLabelService struct {
	inner  labelServicer
	tracer trace.Tracer
}

// NewTracingLabelService returns a TracingLabelService that delegates to inner.
func NewTracingLabelService(inner labelServicer, tracer trace.Tracer) *TracingLabelService {
	return &TracingLabelService{inner: inner, tracer: tracer}
}

func (s *TracingLabelService) Create(ctx context.Context, name string) (label.Label, error) {
	ctx, span := s.tracer.Start(ctx, "tracker.LabelService.Create")
	defer span.End()

	result, err := s.inner.Create(ctx, name)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return label.Label{}, fmt.Errorf("create label: %w", err)
	}
	return result, nil
}

func (s *TracingLabelService) Search(ctx context.Context, name string) ([]label.Label, error) {
	ctx, span := s.tracer.Start(ctx, "tracker.LabelService.Search")
	defer span.End()

	results, err := s.inner.Search(ctx, name)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("search labels: %w", err)
	}
	return results, nil
}
