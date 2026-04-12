package comment

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Comment is the domain entity for a comment on an issue.
type Comment struct {
	ID        uuid.UUID
	Body      string
	AuthorID  uuid.UUID
	IssueID   uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Comments is the paginated result of a List operation.
type Comments struct {
	Items []Comment
}

// ListCommentQuery holds all parameters for listing comments.
type ListCommentQuery struct {
	Cursor *string
	Limit  *int
}

const defaultLimit = 20

// NewListCommentQuery returns a ListCommentQuery with defaults applied.
func NewListCommentQuery(cursor *string, limit *int) ListCommentQuery {
	if limit == nil {
		l := defaultLimit
		limit = &l
	}
	return ListCommentQuery{Cursor: cursor, Limit: limit}
}

// New creates a new Comment with the given fields.
// Returns ErrInvalidComment if any required field is empty or nil.
func New(id uuid.UUID, body string, authorID uuid.UUID, issueID uuid.UUID) (Comment, error) {
	if id == uuid.Nil {
		return Comment{}, ErrInvalidComment
	}
	if body == "" {
		return Comment{}, ErrInvalidComment
	}
	if authorID == uuid.Nil {
		return Comment{}, ErrInvalidComment
	}
	if issueID == uuid.Nil {
		return Comment{}, ErrInvalidComment
	}
	now := time.Now()
	return Comment{
		ID:        id,
		Body:      body,
		AuthorID:  authorID,
		IssueID:   issueID,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Repository defines the persistence operations for comments.
type Repository interface {
	Create(ctx context.Context, comment Comment) (Comment, error)
	Get(ctx context.Context, issueID uuid.UUID) ([]Comment, error)
}

var (
	ErrIssueNotFound  = errors.New("issue not found")
	ErrInvalidComment = errors.New("invalid comment")
)
