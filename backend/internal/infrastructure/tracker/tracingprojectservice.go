package tracker

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	projectdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/project"
)

// projectServicer mirrors api.ProjectService without importing the api package,
// avoiding a layering violation (infrastructure must not depend on api).
type projectServicer interface {
	Create(ctx context.Context, command projectdomain.CreateProjectCommand) (projectdomain.Project, error)
	List(ctx context.Context, query projectdomain.ListProjectQuery) (projectdomain.Projects, error)
}

// TracingProjectService wraps a projectServicer and adds an OTel child span to each operation.
type TracingProjectService struct {
	inner  projectServicer
	tracer trace.Tracer
}

// NewTracingProjectService returns a TracingProjectService that delegates to inner.
func NewTracingProjectService(inner projectServicer, tracer trace.Tracer) *TracingProjectService {
	return &TracingProjectService{inner: inner, tracer: tracer}
}

func (s *TracingProjectService) Create(ctx context.Context, command projectdomain.CreateProjectCommand) (projectdomain.Project, error) {
	ctx, span := s.tracer.Start(ctx, "tracker.ProjectService.Create")
	defer span.End()

	result, err := s.inner.Create(ctx, command)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return projectdomain.Project{}, fmt.Errorf("create project: %w", err)
	}
	return result, nil
}

func (s *TracingProjectService) List(ctx context.Context, query projectdomain.ListProjectQuery) (projectdomain.Projects, error) {
	ctx, span := s.tracer.Start(ctx, "tracker.ProjectService.List")
	defer span.End()

	projects, err := s.inner.List(ctx, query)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return projectdomain.Projects{}, fmt.Errorf("list projects: %w", err)
	}
	return projects, nil
}
