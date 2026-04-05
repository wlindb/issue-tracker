package tracker

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	workspacedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/workspace"
)

var _ workspacedomain.WorkspaceRepository = (*TracingWorkspaceRepository)(nil)

// TracingWorkspaceRepository wraps a WorkspaceRepository and adds an OTel child span to each operation.
type TracingWorkspaceRepository struct {
	inner  workspacedomain.WorkspaceRepository
	tracer trace.Tracer
}

// NewTracingWorkspaceRepository returns a TracingWorkspaceRepository that delegates to inner.
func NewTracingWorkspaceRepository(inner workspacedomain.WorkspaceRepository, tracer trace.Tracer) *TracingWorkspaceRepository {
	return &TracingWorkspaceRepository{inner: inner, tracer: tracer}
}

func (r *TracingWorkspaceRepository) Create(ctx context.Context, workspace workspacedomain.Workspace) (workspacedomain.Workspace, error) {
	ctx, span := r.tracer.Start(ctx, "tracker.WorkspaceRepository.Create")
	defer span.End()

	created, err := r.inner.Create(ctx, workspace)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return workspacedomain.Workspace{}, fmt.Errorf("create workspace: %w", err)
	}
	return created, nil
}

func (r *TracingWorkspaceRepository) Get(ctx context.Context, id uuid.UUID) (*workspacedomain.Workspace, error) {
	ctx, span := r.tracer.Start(ctx, "tracker.WorkspaceRepository.Get")
	defer span.End()

	workspace, err := r.inner.Get(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("get workspace: %w", err)
	}
	return workspace, nil
}

func (r *TracingWorkspaceRepository) List(ctx context.Context, userID uuid.UUID) ([]workspacedomain.Workspace, error) {
	ctx, span := r.tracer.Start(ctx, "tracker.WorkspaceRepository.List")
	defer span.End()

	workspaces, err := r.inner.List(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("list workspaces: %w", err)
	}
	return workspaces, nil
}

func (r *TracingWorkspaceRepository) IsMember(ctx context.Context, workspaceID uuid.UUID, userID uuid.UUID) (bool, error) {
	ctx, span := r.tracer.Start(ctx, "tracker.WorkspaceRepository.IsMember")
	defer span.End()

	member, err := r.inner.IsMember(ctx, workspaceID, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return false, fmt.Errorf("check workspace membership: %w", err)
	}
	return member, nil
}
