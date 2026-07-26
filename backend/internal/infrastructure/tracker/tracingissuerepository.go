package tracker

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	issuedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
	"github.com/wlindb/issue-tracker/internal/domain/tracker/label"
)

var _ issuedomain.IssueRepository = (*TracingIssueRepository)(nil)

// issueRepositoryWithGetByIDs is satisfied by repositories that additionally support GetByIDs,
// a method outside issuedomain.IssueRepository but required by the application layer.
type issueRepositoryWithGetByIDs interface {
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]issuedomain.Issue, error)
}

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

func (r *TracingIssueRepository) GetIssue(ctx context.Context, id uuid.UUID) (issuedomain.Issue, error) {
	ctx, span := r.tracer.Start(ctx, "tracker.IssueRepository.GetIssue")
	defer span.End()

	result, err := r.inner.GetIssue(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return issuedomain.Issue{}, fmt.Errorf("get issue: %w", err)
	}
	return result, nil
}

// GetByIDs delegates to inner if it supports GetByIDs, adding an OTel child span.
func (r *TracingIssueRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]issuedomain.Issue, error) {
	ctx, span := r.tracer.Start(ctx, "tracker.IssueRepository.GetByIDs")
	defer span.End()

	repository, ok := r.inner.(issueRepositoryWithGetByIDs)
	if !ok {
		err := fmt.Errorf("get issues by ids: inner repository does not implement GetByIDs")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	result, err := repository.GetByIDs(ctx, ids)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("get issues by ids: %w", err)
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

func (r *TracingIssueRepository) AddLabel(ctx context.Context, issueID uuid.UUID, label label.Label) error {
	ctx, span := r.tracer.Start(ctx, "tracker.IssueRepository.AddLabel")
	defer span.End()

	if err := r.inner.AddLabel(ctx, issueID, label); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("add label: %w", err)
	}
	return nil
}
