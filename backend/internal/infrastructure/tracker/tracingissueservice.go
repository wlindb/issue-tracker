package tracker

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	issuedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
)

// issueServicer mirrors api.IssueService without importing the api package,
// avoiding a layering violation (infrastructure must not depend on api).
type issueServicer interface {
	ListIssues(ctx context.Context, projectID uuid.UUID, query issuedomain.ListIssueQuery) (issuedomain.IssuePage, error)
	CreateIssue(ctx context.Context, command issuedomain.CreateIssueCommand) (*issuedomain.Issue, error)
	UpdateIssueDescription(ctx context.Context, issueID uuid.UUID, description *string) (*issuedomain.Issue, error)
}

// TracingIssueService wraps an issueServicer and adds an OTel child span to each operation.
type TracingIssueService struct {
	inner  issueServicer
	tracer trace.Tracer
}

// NewTracingIssueService returns a TracingIssueService that delegates to inner.
func NewTracingIssueService(inner issueServicer, tracer trace.Tracer) *TracingIssueService {
	return &TracingIssueService{inner: inner, tracer: tracer}
}

func (s *TracingIssueService) ListIssues(ctx context.Context, projectID uuid.UUID, query issuedomain.ListIssueQuery) (issuedomain.IssuePage, error) {
	ctx, span := s.tracer.Start(ctx, "tracker.IssueService.ListIssues")
	defer span.End()

	page, err := s.inner.ListIssues(ctx, projectID, query)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return issuedomain.IssuePage{}, fmt.Errorf("list issues: %w", err)
	}
	return page, nil
}

func (s *TracingIssueService) CreateIssue(ctx context.Context, command issuedomain.CreateIssueCommand) (*issuedomain.Issue, error) {
	ctx, span := s.tracer.Start(ctx, "tracker.IssueService.CreateIssue")
	defer span.End()

	issue, err := s.inner.CreateIssue(ctx, command)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("create issue: %w", err)
	}
	return issue, nil
}

func (s *TracingIssueService) UpdateIssueDescription(ctx context.Context, issueID uuid.UUID, description *string) (*issuedomain.Issue, error) {
	ctx, span := s.tracer.Start(ctx, "tracker.IssueService.UpdateIssueDescription")
	defer span.End()

	issue, err := s.inner.UpdateIssueDescription(ctx, issueID, description)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("update issue description: %w", err)
	}
	return issue, nil
}
