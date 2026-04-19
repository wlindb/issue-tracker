package tracker

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	commentdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/comment"
)

var _ commentdomain.Repository = (*TracingCommentRepository)(nil)

// TracingCommentRepository wraps a Repository and adds an OTel child span to each operation.
type TracingCommentRepository struct {
	inner  commentdomain.Repository
	tracer trace.Tracer
}

// NewTracingCommentRepository returns a TracingCommentRepository that delegates to inner.
func NewTracingCommentRepository(inner commentdomain.Repository, tracer trace.Tracer) *TracingCommentRepository {
	return &TracingCommentRepository{inner: inner, tracer: tracer}
}

func (r *TracingCommentRepository) Create(ctx context.Context, comment commentdomain.Comment) (commentdomain.Comment, error) {
	ctx, span := r.tracer.Start(ctx, "tracker.CommentRepository.Create")
	defer span.End()

	result, err := r.inner.Create(ctx, comment)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return commentdomain.Comment{}, fmt.Errorf("create comment: %w", err)
	}
	return result, nil
}

func (r *TracingCommentRepository) Get(ctx context.Context, issueID uuid.UUID) ([]commentdomain.Comment, error) {
	ctx, span := r.tracer.Start(ctx, "tracker.CommentRepository.Get")
	defer span.End()

	comments, err := r.inner.Get(ctx, issueID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("get comments: %w", err)
	}
	return comments, nil
}
