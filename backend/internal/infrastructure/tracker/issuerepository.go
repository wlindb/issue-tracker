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

// Compile-time: *IssueRepository must satisfy domain interface.
var _ issuedomain.IssueRepository = (*IssueRepository)(nil)

// IssueRepository is a PostgreSQL-backed implementation of issuedomain.IssueRepository.
type IssueRepository struct {
	db trackerdb.DBTX
}

// NewIssueRepository returns an IssueRepository backed by the given pool.
func NewIssueRepository(db trackerdb.DBTX) *IssueRepository {
	return &IssueRepository{db: db}
}

func (r *IssueRepository) Tx(ctx context.Context, fn func(issuedomain.IssueRepository) error) error {
	pool, ok := r.db.(*pgxpool.Pool)
	if !ok {
		return errors.New("cannot start tx: already inside a transaction")
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if err := fn(NewIssueRepository(tx)); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

// CreateIssue inserts a new issue row and returns the domain model.
func (r *IssueRepository) CreateIssue(ctx context.Context, issue issuedomain.Issue) (issuedomain.Issue, error) {
	queries := trackerdb.New(r.db)

	dbIssue, err := queries.CreateIssue(ctx, createIssueParamsFromDomain(issue))
	if err != nil {
		return issuedomain.Issue{}, fmt.Errorf("create issue: %w", err)
	}

	if err = queries.CreateManyIssueLabels(ctx, createManyIssueLabelsParamsFromDomain(issue)); err != nil {
		return issuedomain.Issue{}, fmt.Errorf("create many issue labels: %w", err)
	}

	return issueToDomain(dbIssue, []issuedomain.Label{}), nil
}

// GetIssue retrieves a single issue by its ID, or ErrIssueNotFound if it does not exist.
func (r *IssueRepository) GetIssue(ctx context.Context, id uuid.UUID) (issuedomain.Issue, error) {
	q := trackerdb.New(r.db)
	row, err := q.GetIssue(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return issuedomain.Issue{}, fmt.Errorf("get issue: %w", issuedomain.ErrIssueNotFound)
		}
		return issuedomain.Issue{}, fmt.Errorf("get issue: %w", err)
	}
	return issueToDomain(row, []issuedomain.Label{}), nil
}

// Update persists mutable fields of an existing issue and returns the updated domain model.
// Returns ErrUpdateConflict if the issue was modified since it was read (optimistic locking).
func (r *IssueRepository) Update(ctx context.Context, issue issuedomain.Issue) (issuedomain.Issue, error) {
	q := trackerdb.New(r.db)
	row, err := q.UpdateIssue(ctx, updateIssueParamsFromDomain(issue))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return issuedomain.Issue{}, fmt.Errorf("update issue: %w", issuedomain.ErrUpdateConflict)
		}
		return issuedomain.Issue{}, fmt.Errorf("update issue: %w", err)
	}
	return issueToDomain(row, []issuedomain.Label{}), nil
}

// ListIssues returns a filtered list of issues for the given project.
func (r *IssueRepository) ListIssues(
	ctx context.Context,
	projectID uuid.UUID,
	query issuedomain.ListIssueQuery,
) (issuedomain.IssuePage, error) {
	q := trackerdb.New(r.db)
	rows, err := q.ListIssues(ctx, listIssuesParamsFromDomain(projectID, query))
	if err != nil {
		return issuedomain.IssuePage{}, fmt.Errorf("list issues: %w", err)
	}

	return issuedomain.IssuePage{
		Items: issuesToDomain(rows),
	}, nil
}
