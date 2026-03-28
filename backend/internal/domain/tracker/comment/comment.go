package comment

import (
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

var (
	ErrIssueNotFound = errors.New("issue not found")
)
