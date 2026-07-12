package tracker

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	userdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/user"
)

var _ userdomain.UserRepository = (*TracingUserRepository)(nil)

// TracingUserRepository wraps a UserRepository and adds an OTel child span to each operation.
type TracingUserRepository struct {
	inner  userdomain.UserRepository
	tracer trace.Tracer
}

// NewTracingUserRepository returns a TracingUserRepository that delegates to inner.
func NewTracingUserRepository(inner userdomain.UserRepository, tracer trace.Tracer) *TracingUserRepository {
	return &TracingUserRepository{inner: inner, tracer: tracer}
}

func (r *TracingUserRepository) Upsert(ctx context.Context, user userdomain.User) (userdomain.User, error) {
	ctx, span := r.tracer.Start(ctx, "tracker.UserRepository.Upsert")
	defer span.End()

	upserted, err := r.inner.Upsert(ctx, user)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return userdomain.User{}, fmt.Errorf("upsert user: %w", err)
	}
	return upserted, nil
}
