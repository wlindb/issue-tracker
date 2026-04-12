package tracker

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	projectdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/project"
)

var _ projectdomain.ProjectRepository = (*TracingProjectRepository)(nil)

// TracingProjectRepository wraps a ProjectRepository and adds an OTel child span to each operation.
type TracingProjectRepository struct {
	inner  projectdomain.ProjectRepository
	tracer trace.Tracer
}

// NewTracingProjectRepository returns a TracingProjectRepository that delegates to inner.
func NewTracingProjectRepository(inner projectdomain.ProjectRepository, tracer trace.Tracer) *TracingProjectRepository {
	return &TracingProjectRepository{inner: inner, tracer: tracer}
}

func (r *TracingProjectRepository) Create(ctx context.Context, project projectdomain.Project) (projectdomain.Project, error) {
	ctx, span := r.tracer.Start(ctx, "tracker.ProjectRepository.Create")
	defer span.End()

	result, err := r.inner.Create(ctx, project)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return projectdomain.Project{}, fmt.Errorf("create project: %w", err)
	}
	return result, nil
}

func (r *TracingProjectRepository) List(ctx context.Context, query projectdomain.ListProjectQuery) (projectdomain.Projects, error) {
	ctx, span := r.tracer.Start(ctx, "tracker.ProjectRepository.List")
	defer span.End()

	projects, err := r.inner.List(ctx, query)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return projectdomain.Projects{}, fmt.Errorf("list projects: %w", err)
	}
	return projects, nil
}

func (r *TracingProjectRepository) Get(ctx context.Context, id uuid.UUID) (projectdomain.Project, error) {
	ctx, span := r.tracer.Start(ctx, "tracker.ProjectRepository.Get")
	defer span.End()

	project, err := r.inner.Get(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return projectdomain.Project{}, fmt.Errorf("get project: %w", err)
	}
	return project, nil
}
