package tracker

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	workspacedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/workspace"
	trackerdb "github.com/wlindb/issue-tracker/internal/infrastructure/tracker/generated"
)

var _ workspacedomain.WorkspaceRepository = (*WorkspaceRepository)(nil)

// WorkspaceRepository is a PostgreSQL-backed implementation of workspacedomain.WorkspaceRepository.
type WorkspaceRepository struct {
	pool    *pgxpool.Pool
	queries *trackerdb.Queries
}

// NewWorkspaceRepository returns a WorkspaceRepository backed by the given pool.
func NewWorkspaceRepository(pool *pgxpool.Pool) *WorkspaceRepository {
	return &WorkspaceRepository{
		pool:    pool,
		queries: trackerdb.New(pool),
	}
}

// Create inserts a workspace row and its owner as a member atomically inside a transaction.
func (r *WorkspaceRepository) Create(ctx context.Context, workspace workspacedomain.Workspace) (workspacedomain.Workspace, error) {
	transaction, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return workspacedomain.Workspace{}, fmt.Errorf("create workspace: begin transaction: %w", err)
	}
	defer func() { _ = transaction.Rollback(ctx) }()

	queries := trackerdb.New(transaction)

	row, err := queries.CreateWorkspace(ctx, trackerdb.CreateWorkspaceParams{
		ID:      workspace.ID,
		Name:    workspace.Name,
		OwnerID: workspace.OwnerID,
	})
	if err != nil {
		return workspacedomain.Workspace{}, fmt.Errorf("create workspace: %w", err)
	}

	if err = queries.InsertWorkspaceMember(ctx, trackerdb.InsertWorkspaceMemberParams{
		WorkspaceID: workspace.ID,
		UserID:      workspace.OwnerID,
	}); err != nil {
		return workspacedomain.Workspace{}, fmt.Errorf("create workspace: insert member: %w", err)
	}

	if err = transaction.Commit(ctx); err != nil {
		return workspacedomain.Workspace{}, fmt.Errorf("create workspace: commit: %w", err)
	}

	return workspaceToDomain(row), nil
}

// Get returns a single workspace by ID, or ErrWorkspaceNotFound if it does not exist.
func (r *WorkspaceRepository) Get(ctx context.Context, id uuid.UUID) (*workspacedomain.Workspace, error) {
	row, err := r.queries.GetWorkspace(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("get workspace: %w", workspacedomain.ErrWorkspaceNotFound)
		}
		return nil, fmt.Errorf("get workspace: %w", err)
	}
	w := workspaceToDomain(row)
	return &w, nil
}

// List returns all workspaces the given user is a member of.
func (r *WorkspaceRepository) List(ctx context.Context, userID uuid.UUID) ([]workspacedomain.Workspace, error) {
	rows, err := r.queries.ListWorkspacesForUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list workspaces: %w", err)
	}
	return workspacesToDomain(rows), nil
}

// IsMember reports whether userID is a member of workspaceID.
func (r *WorkspaceRepository) IsMember(ctx context.Context, workspaceID uuid.UUID, userID uuid.UUID) (bool, error) {
	member, err := r.queries.IsMember(ctx, trackerdb.IsMemberParams{
		WorkspaceID: workspaceID,
		UserID:      userID,
	})
	if err != nil {
		return false, fmt.Errorf("check workspace membership: %w", err)
	}
	return member, nil
}
