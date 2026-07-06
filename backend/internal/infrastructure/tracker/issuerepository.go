package tracker

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	issuedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
	"github.com/wlindb/issue-tracker/internal/domain/tracker/label"
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

	return issueToDomain(dbIssue, []label.Label{}), nil
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
	return issueToDomain(row, []label.Label{}), nil
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
	return issueToDomain(row, []label.Label{}), nil
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

// AddLabel attaches the given label to the issue, idempotently (attaching an
// already-present label is a success, not an error). Returns label.ErrLabelNotFound
// if the label does not exist, or issuedomain.ErrIssueNotFound if the issue does not exist.
func (r *IssueRepository) AddLabel(ctx context.Context, issueID uuid.UUID, l label.Label) error {
	queries := trackerdb.New(r.db)
	err := queries.AddIssueLabel(ctx, trackerdb.AddIssueLabelParams{
		IssueID: issueID,
		LabelID: l.ID,
	})
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
		switch pgErr.ConstraintName {
		case "issue_labels_label_id_fkey", "issue_labels_workspace_matches_label":
			return fmt.Errorf("add label: %w", label.ErrLabelNotFound)
		case "issue_labels_issue_id_fkey", "issue_labels_workspace_matches_issue":
			return fmt.Errorf("add label: %w", issuedomain.ErrIssueNotFound)
		}
	}
	return fmt.Errorf("add label: %w", err)
}
