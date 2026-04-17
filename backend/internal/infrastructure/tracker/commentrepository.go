package tracker

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	commentdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/comment"
	trackerdb "github.com/wlindb/issue-tracker/internal/infrastructure/tracker/generated"
)

// commentQuerier defines the query methods used by CommentRepository.
type commentQuerier interface {
	CreateComment(ctx context.Context, arg trackerdb.CreateCommentParams) (trackerdb.Comment, error)
	ListCommentsByIssue(ctx context.Context, issueID uuid.UUID) ([]trackerdb.Comment, error)
}

// Compile-time: *CommentRepository must satisfy domain interface.
var _ commentdomain.Repository = (*CommentRepository)(nil)

// CommentRepository is a PostgreSQL-backed implementation of commentdomain.Repository.
type CommentRepository struct {
	queries commentQuerier
}

// NewCommentRepository returns a CommentRepository backed by the given pool.
func NewCommentRepository(pool *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{queries: trackerdb.New(pool)}
}

// Create inserts a new comment row and returns the domain model.
func (r *CommentRepository) Create(ctx context.Context, comment commentdomain.Comment) (commentdomain.Comment, error) {
	row, err := r.queries.CreateComment(ctx, createCommentParamsFromDomain(comment))
	if err != nil {
		return commentdomain.Comment{}, fmt.Errorf("create comment: %w", err)
	}
	return commentToDomain(row), nil
}

// Get retrieves all comments for the given issue ID.
func (r *CommentRepository) Get(ctx context.Context, issueID uuid.UUID) ([]commentdomain.Comment, error) {
	rows, err := r.queries.ListCommentsByIssue(ctx, issueID)
	if err != nil {
		return nil, fmt.Errorf("get comments: %w", err)
	}
	return commentsToDomain(rows), nil
}
