//go:build !integration

package comment_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wlindb/issue-tracker/internal/domain/tracker/comment"
)

func Test_New_NilID_ReturnsError(t *testing.T) {
	actual, err := comment.New(uuid.Nil, "body", uuid.New(), uuid.New())

	require.Error(t, err)
	assert.ErrorIs(t, err, comment.ErrInvalidComment)
	assert.Zero(t, actual)
}

func Test_New_EmptyBody_ReturnsError(t *testing.T) {
	actual, err := comment.New(uuid.New(), "", uuid.New(), uuid.New())

	require.Error(t, err)
	assert.ErrorIs(t, err, comment.ErrInvalidComment)
	assert.Zero(t, actual)
}

func Test_New_NilAuthorID_ReturnsError(t *testing.T) {
	actual, err := comment.New(uuid.New(), "body", uuid.Nil, uuid.New())

	require.Error(t, err)
	assert.ErrorIs(t, err, comment.ErrInvalidComment)
	assert.Zero(t, actual)
}

func Test_New_NilIssueID_ReturnsError(t *testing.T) {
	actual, err := comment.New(uuid.New(), "body", uuid.New(), uuid.Nil)

	require.Error(t, err)
	assert.ErrorIs(t, err, comment.ErrInvalidComment)
	assert.Zero(t, actual)
}

func Test_New_ValidArguments_ReturnsComment(t *testing.T) {
	id := uuid.New()
	authorID := uuid.New()
	issueID := uuid.New()

	actual, err := comment.New(id, "my comment", authorID, issueID)

	require.NoError(t, err)
	assert.Equal(t, id, actual.ID)
	assert.Equal(t, "my comment", actual.Body)
	assert.Equal(t, authorID, actual.AuthorID)
	assert.Equal(t, issueID, actual.IssueID)
	assert.False(t, actual.CreatedAt.IsZero())
	assert.False(t, actual.UpdatedAt.IsZero())
}

func Test_NewListCommentQuery_NilLimit_UsesDefaultLimit(t *testing.T) {
	actual := comment.NewListCommentQuery(nil, nil)

	require.NotNil(t, actual.Limit)
	assert.Equal(t, 20, *actual.Limit)
}

func Test_NewListCommentQuery_ExplicitLimit_UsesProvidedLimit(t *testing.T) {
	limit := 5

	actual := comment.NewListCommentQuery(nil, &limit)

	require.NotNil(t, actual.Limit)
	assert.Equal(t, 5, *actual.Limit)
}

func Test_NewListCommentQuery_NilCursor_CursorIsNil(t *testing.T) {
	actual := comment.NewListCommentQuery(nil, nil)

	assert.Nil(t, actual.Cursor)
}
