package tracker

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	issuedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
)

var _ issuedomain.IssueRepository = (*TracingIssueRepository)(nil)

// TracingIssueRepository wraps an IssueRepository and adds an OTel child span to each operation.
type TracingIssueRepository struct {
	inner  issuedomain.IssueRepository
	tracer trace.Tracer
}

// NewTracingIssueRepository returns a TracingIssueRepository that delegates to inner.
func NewTracingIssueRepository(inner issuedomain.IssueRepository, tracer trace.Tracer) *TracingIssueRepository {
	return &TracingIssueRepository{inner: inner, tracer: tracer}
}

func (r *TracingIssueRepository) ListIssues(ctx context.Context, projectID uuid.UUID, query issuedomain.ListIssueQuery) (issuedomain.IssuePage, error) {
	ctx, span := r.tracer.Start(ctx, "tracker.IssueRepository.ListIssues")
	defer span.End()

	page, err := r.inner.ListIssues(ctx, projectID, query)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return issuedomain.IssuePage{}, fmt.Errorf("list issues: %w", err)
	}
	return page, nil
}

func (r *TracingIssueRepository) CreateIssue(ctx context.Context, issue issuedomain.Issue) (issuedomain.Issue, error) {
	ctx, span := r.tracer.Start(ctx, "tracker.IssueRepository.CreateIssue")
	defer span.End()

	result, err := r.inner.CreateIssue(ctx, issue)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return issuedomain.Issue{}, fmt.Errorf("create issue: %w", err)
	}
	return result, nil
}

func (r *TracingIssueRepository) Update(ctx context.Context, issue issuedomain.Issue) (issuedomain.Issue, error) {
	ctx, span := r.tracer.Start(ctx, "tracker.IssueRepository.Update")
	defer span.End()

	result, err := r.inner.Update(ctx, issue)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return issuedomain.Issue{}, fmt.Errorf("update issue: %w", err)
	}
	return result, nil
}
