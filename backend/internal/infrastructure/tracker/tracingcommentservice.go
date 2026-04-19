package tracker

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	commentdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/comment"
)

// commentServicer mirrors api.CommentService without importing the api package,
// avoiding a layering violation (infrastructure must not depend on api).
type commentServicer interface {
	List(ctx context.Context, issueID uuid.UUID, query commentdomain.ListCommentQuery) (commentdomain.Comments, error)
	Create(ctx context.Context, comment commentdomain.Comment) (*commentdomain.Comment, error)
}

// TracingCommentService wraps a commentServicer and adds an OTel child span to each operation.
type TracingCommentService struct {
	inner  commentServicer
	tracer trace.Tracer
}

// NewTracingCommentService returns a TracingCommentService that delegates to inner.
func NewTracingCommentService(inner commentServicer, tracer trace.Tracer) *TracingCommentService {
	return &TracingCommentService{inner: inner, tracer: tracer}
}

func (s *TracingCommentService) List(ctx context.Context, issueID uuid.UUID, query commentdomain.ListCommentQuery) (commentdomain.Comments, error) {
	ctx, span := s.tracer.Start(ctx, "tracker.CommentService.List")
	defer span.End()

	comments, err := s.inner.List(ctx, issueID, query)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return commentdomain.Comments{}, fmt.Errorf("list comments: %w", err)
	}
	return comments, nil
}

func (s *TracingCommentService) Create(ctx context.Context, comment commentdomain.Comment) (*commentdomain.Comment, error) {
	ctx, span := s.tracer.Start(ctx, "tracker.CommentService.Create")
	defer span.End()

	result, err := s.inner.Create(ctx, comment)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("create comment: %w", err)
	}
	return result, nil
}
