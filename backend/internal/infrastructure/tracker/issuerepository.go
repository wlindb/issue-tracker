package tracker

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	issuedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
	trackerdb "github.com/wlindb/issue-tracker/internal/infrastructure/tracker/generated"
)

// issueQuerier defines the query methods used by IssueRepository.
type issueQuerier interface {
	GetIssue(ctx context.Context, id uuid.UUID) (trackerdb.Issue, error)
	CreateIssue(ctx context.Context, arg trackerdb.CreateIssueParams) (trackerdb.Issue, error)
	ListIssues(ctx context.Context, arg trackerdb.ListIssuesParams) ([]trackerdb.Issue, error)
	UpdateIssue(ctx context.Context, arg trackerdb.UpdateIssueParams) (trackerdb.Issue, error)
}

// Compile-time: *IssueRepository must satisfy domain interface.
var _ issuedomain.IssueRepository = (*IssueRepository)(nil)

// IssueRepository is a PostgreSQL-backed implementation of issuedomain.IssueRepository.
type IssueRepository struct {
	queries issueQuerier
}

// NewIssueRepository returns an IssueRepository backed by the given pool.
func NewIssueRepository(pool *pgxpool.Pool) *IssueRepository {
	return &IssueRepository{queries: trackerdb.New(pool)}
}

// CreateIssue inserts a new issue row and returns the domain model.
func (r *IssueRepository) CreateIssue(ctx context.Context, issue issuedomain.Issue) (issuedomain.Issue, error) {
	row, err := r.queries.CreateIssue(ctx, createIssueParamsFromDomain(issue))
	if err != nil {
		return issuedomain.Issue{}, fmt.Errorf("create issue: %w", err)
	}
	return issueToDomain(row), nil
}

// GetIssue retrieves a single issue by its ID, or ErrIssueNotFound if it does not exist.
func (r *IssueRepository) GetIssue(ctx context.Context, id uuid.UUID) (issuedomain.Issue, error) {
	row, err := r.queries.GetIssue(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return issuedomain.Issue{}, fmt.Errorf("get issue: %w", issuedomain.ErrIssueNotFound)
		}
		return issuedomain.Issue{}, fmt.Errorf("get issue: %w", err)
	}
	return issueToDomain(row), nil
}

// Update persists mutable fields of an existing issue and returns the updated domain model.
// Returns ErrUpdateConflict if the issue was modified since it was read (optimistic locking).
func (r *IssueRepository) Update(ctx context.Context, issue issuedomain.Issue) (issuedomain.Issue, error) {
	row, err := r.queries.UpdateIssue(ctx, updateIssueParamsFromDomain(issue))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return issuedomain.Issue{}, fmt.Errorf("update issue: %w", issuedomain.ErrUpdateConflict)
		}
		return issuedomain.Issue{}, fmt.Errorf("update issue: %w", err)
	}
	return issueToDomain(row), nil
}

// ListIssues returns a filtered list of issues for the given project.
func (r *IssueRepository) ListIssues(
	ctx context.Context,
	projectID uuid.UUID,
	query issuedomain.ListIssueQuery,
) (issuedomain.IssuePage, error) {
	rows, err := r.queries.ListIssues(ctx, listIssuesParamsFromDomain(projectID, query))
	if err != nil {
		return issuedomain.IssuePage{}, fmt.Errorf("list issues: %w", err)
	}

	return issuedomain.IssuePage{
		Items: issuesToDomain(rows),
	}, nil
}
