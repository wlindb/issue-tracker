//go:build !integration

package search_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/wlindb/issue-tracker/internal/application/search"
	"github.com/wlindb/issue-tracker/internal/domain/search/issuedocument"
	"github.com/wlindb/issue-tracker/internal/pkg/domain/event"
	trackermodel "github.com/wlindb/issue-tracker/internal/pkg/tracker/model"
)

type fakeSubscriber[T any] struct {
	handler event.Subscriber[T]
}

func (s *fakeSubscriber[T]) Subscribe(handler event.Subscriber[T]) error {
	s.handler = handler
	return nil
}

type failingSubscriber[T any] struct {
	err error
}

func (s *failingSubscriber[T]) Subscribe(_ event.Subscriber[T]) error {
	return s.err
}

type mockIssueDocumentRepository struct {
	mock.Mock
}

func (m *mockIssueDocumentRepository) Create(ctx context.Context, document issuedocument.IssueDocument) (issuedocument.IssueDocument, error) {
	args := m.Called(ctx, document)
	if d, ok := args.Get(0).(issuedocument.IssueDocument); ok {
		return d, args.Error(1)
	}
	return issuedocument.IssueDocument{}, args.Error(1)
}

func (m *mockIssueDocumentRepository) Update(ctx context.Context, document issuedocument.IssueDocument) (issuedocument.IssueDocument, error) {
	args := m.Called(ctx, document)
	if d, ok := args.Get(0).(issuedocument.IssueDocument); ok {
		return d, args.Error(1)
	}
	return issuedocument.IssueDocument{}, args.Error(1)
}

func (m *mockIssueDocumentRepository) Find(ctx context.Context, description string) (issuedocument.IssueDocuments, error) {
	args := m.Called(ctx, description)
	if d, ok := args.Get(0).(issuedocument.IssueDocuments); ok {
		return d, args.Error(1)
	}
	return issuedocument.IssueDocuments{}, args.Error(1)
}

func Test_NewSearchEventHandler_WithIssueCreated_SubscribesHandleIssueCreated(t *testing.T) {
	subscriber := &fakeSubscriber[trackermodel.IssueCreatedEvent]{}
	repository := &mockIssueDocumentRepository{}
	repository.On("Create", mock.Anything, mock.Anything).Return(issuedocument.IssueDocument{}, nil)

	_, err := search.NewSearchEventHandler(repository, search.WithIssueCreated(subscriber))
	require.NoError(t, err)
	require.NotNil(t, subscriber.handler)

	err = subscriber.handler(context.Background(), trackermodel.IssueCreatedEvent{
		OccurredAt: time.Now(),
		Payload: trackermodel.Issue{
			ID:         uuid.New(),
			ReporterID: uuid.New(),
			Title:      "test issue",
		},
	})
	require.NoError(t, err)
}

func Test_NewSearchEventHandler_WithIssueTitleUpdated_SubscribesHandleIssueTitleUpdated(t *testing.T) {
	subscriber := &fakeSubscriber[trackermodel.IssueTitleUpdatedEvent]{}
	repository := &mockIssueDocumentRepository{}

	_, err := search.NewSearchEventHandler(repository, search.WithIssueTitleUpdated(subscriber))
	require.NoError(t, err)
	require.NotNil(t, subscriber.handler)

	err = subscriber.handler(context.Background(), trackermodel.IssueTitleUpdatedEvent{
		OccurredAt: time.Now(),
		Payload: trackermodel.Issue{
			ID:    uuid.New(),
			Title: "updated title",
		},
	})
	require.NoError(t, err)
}

func Test_NewSearchEventHandler_WithIssueDescriptionUpdated_SubscribesHandleIssueDescriptionUpdated(t *testing.T) {
	subscriber := &fakeSubscriber[trackermodel.IssueDescriptionUpdatedEvent]{}
	repository := &mockIssueDocumentRepository{}

	_, err := search.NewSearchEventHandler(repository, search.WithIssueDescriptionUpdated(subscriber))
	require.NoError(t, err)
	require.NotNil(t, subscriber.handler)

	description := "updated description"
	err = subscriber.handler(context.Background(), trackermodel.IssueDescriptionUpdatedEvent{
		OccurredAt: time.Now(),
		Payload: trackermodel.Issue{
			ID:          uuid.New(),
			Description: &description,
		},
	})
	require.NoError(t, err)
}

func Test_NewSearchEventHandler_SubscribeError_ReturnsError(t *testing.T) {
	expected := errors.New("subscribe failed")
	subscriber := &failingSubscriber[trackermodel.IssueCreatedEvent]{err: expected}
	repository := &mockIssueDocumentRepository{}

	_, err := search.NewSearchEventHandler(repository, search.WithIssueCreated(subscriber))
	require.ErrorIs(t, err, expected)
}

func Test_HandleIssueCreated_ValidEvent_CreatesIssueDocument(t *testing.T) {
	repository := &mockIssueDocumentRepository{}
	occurredAt := time.Now()
	issueID := uuid.New()
	description := "test description"

	expected := issuedocument.NewIssueDocument(issueID, "test issue", description, occurredAt)
	repository.On("Create", mock.Anything, expected).Return(expected, nil)

	handler, err := search.NewSearchEventHandler(repository)
	require.NoError(t, err)

	err = handler.HandleIssueCreated(context.Background(), trackermodel.IssueCreatedEvent{
		OccurredAt: occurredAt,
		Payload: trackermodel.Issue{
			ID:          issueID,
			Title:       "test issue",
			Description: &description,
		},
	})

	require.NoError(t, err)
	repository.AssertExpectations(t)
}

func Test_HandleIssueCreated_NilDescription_CreatesIssueDocumentWithEmptyDescription(t *testing.T) {
	repository := &mockIssueDocumentRepository{}
	occurredAt := time.Now()
	issueID := uuid.New()

	expected := issuedocument.NewIssueDocument(issueID, "test issue", "", occurredAt)
	repository.On("Create", mock.Anything, expected).Return(expected, nil)

	handler, err := search.NewSearchEventHandler(repository)
	require.NoError(t, err)

	err = handler.HandleIssueCreated(context.Background(), trackermodel.IssueCreatedEvent{
		OccurredAt: occurredAt,
		Payload: trackermodel.Issue{
			ID:    issueID,
			Title: "test issue",
		},
	})

	require.NoError(t, err)
	repository.AssertExpectations(t)
}

func Test_HandleIssueCreated_RepositoryError_ReturnsError(t *testing.T) {
	repository := &mockIssueDocumentRepository{}
	expected := errors.New("create failed")
	repository.On("Create", mock.Anything, mock.Anything).Return(issuedocument.IssueDocument{}, expected)

	handler, err := search.NewSearchEventHandler(repository)
	require.NoError(t, err)

	err = handler.HandleIssueCreated(context.Background(), trackermodel.IssueCreatedEvent{
		OccurredAt: time.Now(),
		Payload: trackermodel.Issue{
			ID:    uuid.New(),
			Title: "test issue",
		},
	})

	require.ErrorIs(t, err, expected)
}
