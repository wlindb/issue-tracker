//go:build !integration

package tracker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	commentdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/comment"
	trackerdb "github.com/wlindb/issue-tracker/internal/infrastructure/tracker/generated"
)

type mockCommentQuerier struct {
	mock.Mock
}

func (m *mockCommentQuerier) CreateComment(ctx context.Context, arg trackerdb.CreateCommentParams) (trackerdb.Comment, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(trackerdb.Comment), args.Error(1)
}

func (m *mockCommentQuerier) ListCommentsByIssue(ctx context.Context, issueID uuid.UUID) ([]trackerdb.Comment, error) {
	args := m.Called(ctx, issueID)
	if result, ok := args.Get(0).([]trackerdb.Comment); ok {
		return result, args.Error(1)
	}
	return nil, args.Error(1)
}

// — Create unit tests —

func Test_CreateComment_Success_ReturnsDomainComment(t *testing.T) {
	querier := &mockCommentQuerier{}
	repository := &CommentRepository{queries: querier}

	now := time.Now().UTC()
	domainComment := commentdomain.Comment{
		ID:       uuid.New(),
		Body:     "Test comment",
		AuthorID: uuid.New(),
		IssueID:  uuid.New(),
	}

	returnedRow := trackerdb.Comment{
		ID:        domainComment.ID,
		Body:      domainComment.Body,
		AuthorID:  domainComment.AuthorID,
		IssueID:   domainComment.IssueID,
		CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: now, Valid: true},
	}

	querier.On("CreateComment", mock.Anything, mock.Anything).Return(returnedRow, nil)

	actual, err := repository.Create(context.Background(), domainComment)

	require.NoError(t, err)
	assert.Equal(t, domainComment.ID, actual.ID)
	assert.Equal(t, domainComment.Body, actual.Body)
	assert.Equal(t, domainComment.AuthorID, actual.AuthorID)
	assert.Equal(t, domainComment.IssueID, actual.IssueID)
	querier.AssertExpectations(t)
}

func Test_CreateComment_QueryError_ReturnsWrappedError(t *testing.T) {
	querier := &mockCommentQuerier{}
	repository := &CommentRepository{queries: querier}

	dbErr := errors.New("foreign key violation")
	querier.On("CreateComment", mock.Anything, mock.Anything).Return(trackerdb.Comment{}, dbErr)

	_, err := repository.Create(context.Background(), commentdomain.Comment{
		ID:       uuid.New(),
		Body:     "Test comment",
		AuthorID: uuid.New(),
		IssueID:  uuid.New(),
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.Contains(t, err.Error(), "create comment")
	querier.AssertExpectations(t)
}

// — Get unit tests —

func Test_GetComments_Success_ReturnsDomainComments(t *testing.T) {
	querier := &mockCommentQuerier{}
	repository := &CommentRepository{queries: querier}

	issueID := uuid.New()
	now := time.Now().UTC()

	returnedRows := []trackerdb.Comment{
		{
			ID:        uuid.New(),
			Body:      "First comment",
			AuthorID:  uuid.New(),
			IssueID:   issueID,
			CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
			UpdatedAt: pgtype.Timestamptz{Time: now, Valid: true},
		},
	}

	querier.On("ListCommentsByIssue", mock.Anything, issueID).Return(returnedRows, nil)

	actual, err := repository.Get(context.Background(), issueID)

	require.NoError(t, err)
	require.Len(t, actual, 1)
	assert.Equal(t, "First comment", actual[0].Body)
	assert.Equal(t, issueID, actual[0].IssueID)
	querier.AssertExpectations(t)
}

func Test_GetComments_EmptyResult_ReturnsEmptySlice(t *testing.T) {
	querier := &mockCommentQuerier{}
	repository := &CommentRepository{queries: querier}

	querier.On("ListCommentsByIssue", mock.Anything, mock.Anything).Return([]trackerdb.Comment{}, nil)

	actual, err := repository.Get(context.Background(), uuid.New())

	require.NoError(t, err)
	assert.Empty(t, actual)
	querier.AssertExpectations(t)
}

func Test_GetComments_QueryError_ReturnsWrappedError(t *testing.T) {
	querier := &mockCommentQuerier{}
	repository := &CommentRepository{queries: querier}

	dbErr := errors.New("connection refused")
	querier.On("ListCommentsByIssue", mock.Anything, mock.Anything).Return([]trackerdb.Comment(nil), dbErr)

	_, err := repository.Get(context.Background(), uuid.New())

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.Contains(t, err.Error(), "get comments")
	querier.AssertExpectations(t)
}
