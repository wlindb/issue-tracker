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
	ListLabelsByIssueIDs(ctx context.Context, issueIDs []uuid.UUID) ([]trackerdb.ListLabelsByIssueIDsRow, error)
	GetLabelsByIDs(ctx context.Context, ids []uuid.UUID) ([]trackerdb.GetLabelsByIDsRow, error)
}

// Compile-time: *IssueRepository must satisfy domain interface.
var _ issuedomain.IssueRepository = (*IssueRepository)(nil)

// IssueRepository is a PostgreSQL-backed implementation of issuedomain.IssueRepository.
type IssueRepository struct {
	pool    *pgxpool.Pool
	queries issueQuerier
}

// NewIssueRepository returns an IssueRepository backed by the given pool.
func NewIssueRepository(pool *pgxpool.Pool) *IssueRepository {
	return &IssueRepository{pool: pool, queries: trackerdb.New(pool)}
}

// CreateIssue inserts a new issue row together with its labels in a single transaction.
func (r *IssueRepository) CreateIssue(ctx context.Context, issue issuedomain.Issue) (issuedomain.Issue, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return issuedomain.Issue{}, fmt.Errorf("create issue begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	txQueries := trackerdb.New(tx)

	row, err := txQueries.CreateIssue(ctx, createIssueParamsFromDomain(issue))
	if err != nil {
		return issuedomain.Issue{}, fmt.Errorf("create issue: %w", err)
	}

	for _, label := range issue.Labels {
		if err := txQueries.InsertIssueLabel(ctx, trackerdb.InsertIssueLabelParams{
			IssueID: row.ID,
			LabelID: label.ID,
		}); err != nil {
			return issuedomain.Issue{}, fmt.Errorf("insert issue label: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return issuedomain.Issue{}, fmt.Errorf("create issue commit: %w", err)
	}

	return issueToDomain(row, issue.Labels), nil
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

	labelRows, err := r.queries.ListLabelsByIssueIDs(ctx, []uuid.UUID{id})
	if err != nil {
		return issuedomain.Issue{}, fmt.Errorf("get issue labels: %w", err)
	}

	labelMap := labelsByIssueFromDB(labelRows)
	return issueToDomain(row, labelMap[id]), nil
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
	return issueToDomain(row, issue.Labels), nil
}

// ListIssues returns a filtered list of issues for the given project, with their labels.
func (r *IssueRepository) ListIssues(
	ctx context.Context,
	projectID uuid.UUID,
	query issuedomain.ListIssueQuery,
) (issuedomain.IssuePage, error) {
	rows, err := r.queries.ListIssues(ctx, listIssuesParamsFromDomain(projectID, query))
	if err != nil {
		return issuedomain.IssuePage{}, fmt.Errorf("list issues: %w", err)
	}

	if len(rows) == 0 {
		return issuedomain.IssuePage{Items: []issuedomain.Issue{}}, nil
	}

	issueIDs := make([]uuid.UUID, len(rows))
	for i, row := range rows {
		issueIDs[i] = row.ID
	}

	labelRows, err := r.queries.ListLabelsByIssueIDs(ctx, issueIDs)
	if err != nil {
		return issuedomain.IssuePage{}, fmt.Errorf("list issue labels: %w", err)
	}

	labelsByIssue := labelsByIssueFromDB(labelRows)

	items := make([]issuedomain.Issue, len(rows))
	for i, row := range rows {
		items[i] = issueToDomain(row, labelsByIssue[row.ID])
	}

	return issuedomain.IssuePage{Items: items}, nil
}

// GetLabelsByIDs resolves a slice of label IDs to domain Label values.
func (r *IssueRepository) GetLabelsByIDs(ctx context.Context, ids []uuid.UUID) ([]issuedomain.Label, error) {
	if len(ids) == 0 {
		return []issuedomain.Label{}, nil
	}
	rows, err := r.queries.GetLabelsByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("get labels by IDs: %w", err)
	}
	return labelsFromDB(rows), nil
}
