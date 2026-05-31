//go:build !integration

package comment_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/wlindb/issue-tracker/internal/domain/tracker/comment"
	"github.com/wlindb/issue-tracker/internal/pkg/domain/event"
)

type mockCommentRepository struct {
	mock.Mock
}

func (m *mockCommentRepository) Create(ctx context.Context, c comment.Comment) (comment.Comment, error) {
	args := m.Called(ctx, c)
	if result, ok := args.Get(0).(comment.Comment); ok {
		return result, args.Error(1)
	}
	return comment.Comment{}, args.Error(1)
}

func (m *mockCommentRepository) Get(ctx context.Context, issueID uuid.UUID) ([]comment.Comment, error) {
	args := m.Called(ctx, issueID)
	if result, ok := args.Get(0).([]comment.Comment); ok {
		return result, args.Error(1)
	}
	return nil, args.Error(1)
}

func Test_Create_SuccessfulPersistence_PublishesCommentCreatedEvent(t *testing.T) {
	repository := &mockCommentRepository{}
	service := comment.NewCommentService(repository)

	c := comment.Comment{
		ID:       uuid.New(),
		Body:     "hello",
		AuthorID: uuid.New(),
		IssueID:  uuid.New(),
	}

	var published []comment.CommentCreatedEvent
	ctx := event.WithPublisher[comment.CommentCreatedEvent](context.Background(), func(_ context.Context, e comment.CommentCreatedEvent) error {
		published = append(published, e)
		return nil
	})

	repository.On("Create", mock.Anything, c).Return(c, nil)

	result, err := service.Create(ctx, c)
	require.NoError(t, err)
	require.NotNil(t, result)

	require.Len(t, published, 1)
	assert.Equal(t, c, published[0].Payload)
	assert.False(t, published[0].OccurredAt.IsZero())
	repository.AssertExpectations(t)
}

func Test_Create_PublisherError_StillReturnsComment(t *testing.T) {
	repository := &mockCommentRepository{}
	service := comment.NewCommentService(repository)

	c := comment.Comment{
		ID:       uuid.New(),
		Body:     "hello",
		AuthorID: uuid.New(),
		IssueID:  uuid.New(),
	}

	ctx := event.WithPublisher[comment.CommentCreatedEvent](context.Background(), func(_ context.Context, _ comment.CommentCreatedEvent) error {
		return errors.New("nats down")
	})

	repository.On("Create", mock.Anything, c).Return(c, nil)

	result, err := service.Create(ctx, c)
	require.NoError(t, err)
	assert.Equal(t, &c, result)
	repository.AssertExpectations(t)
}

func Test_Create_RepositoryError_ReturnsError(t *testing.T) {
	repository := &mockCommentRepository{}
	service := comment.NewCommentService(repository)

	c := comment.Comment{
		ID:       uuid.New(),
		Body:     "hello",
		AuthorID: uuid.New(),
		IssueID:  uuid.New(),
	}
	repositoryErr := errors.New("db error")

	repository.On("Create", mock.Anything, c).Return(comment.Comment{}, repositoryErr)

	result, err := service.Create(context.Background(), c)
	require.Error(t, err)
	assert.ErrorIs(t, err, repositoryErr)
	assert.Nil(t, result)
	repository.AssertExpectations(t)
}
