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

type mockIssueDocumentService struct {
	mock.Mock
}

func (m *mockIssueDocumentService) Create(ctx context.Context, id uuid.UUID, title string, description string, updatedAt time.Time) error {
	args := m.Called(ctx, id, title, description, updatedAt)
	return args.Error(0)
}

func (m *mockIssueDocumentService) UpdateTitle(ctx context.Context, id uuid.UUID, title string) error {
	args := m.Called(ctx, id, title)
	return args.Error(0)
}

func (m *mockIssueDocumentService) UpdateDescription(ctx context.Context, id uuid.UUID, description string) error {
	args := m.Called(ctx, id, description)
	return args.Error(0)
}

func Test_NewSearchEventHandler_WithIssueCreated_SubscribesHandleIssueCreated(t *testing.T) {
	subscriber := &fakeSubscriber[trackermodel.IssueCreatedEvent]{}
	service := &mockIssueDocumentService{}
	service.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	_, err := search.NewSearchEventHandler(service, search.WithIssueCreated(subscriber))
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
	service := &mockIssueDocumentService{}
	service.On("UpdateTitle", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	_, err := search.NewSearchEventHandler(service, search.WithIssueTitleUpdated(subscriber))
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
	service := &mockIssueDocumentService{}
	service.On("UpdateDescription", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	_, err := search.NewSearchEventHandler(service, search.WithIssueDescriptionUpdated(subscriber))
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
	service := &mockIssueDocumentService{}

	_, err := search.NewSearchEventHandler(service, search.WithIssueCreated(subscriber))
	require.ErrorIs(t, err, expected)
}

func Test_HandleIssueCreated_ValidEvent_CreatesIssueDocument(t *testing.T) {
	service := &mockIssueDocumentService{}
	occurredAt := time.Now()
	issueID := uuid.New()
	description := "test description"

	service.On("Create", mock.Anything, issueID, "test issue", description, occurredAt).Return(nil)

	handler, err := search.NewSearchEventHandler(service)
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
	service.AssertExpectations(t)
}

func Test_HandleIssueCreated_NilDescription_CreatesIssueDocumentWithEmptyDescription(t *testing.T) {
	service := &mockIssueDocumentService{}
	occurredAt := time.Now()
	issueID := uuid.New()

	service.On("Create", mock.Anything, issueID, "test issue", "", occurredAt).Return(nil)

	handler, err := search.NewSearchEventHandler(service)
	require.NoError(t, err)

	err = handler.HandleIssueCreated(context.Background(), trackermodel.IssueCreatedEvent{
		OccurredAt: occurredAt,
		Payload: trackermodel.Issue{
			ID:    issueID,
			Title: "test issue",
		},
	})

	require.NoError(t, err)
	service.AssertExpectations(t)
}

func Test_HandleIssueCreated_RepositoryError_ReturnsError(t *testing.T) {
	service := &mockIssueDocumentService{}
	expected := errors.New("create failed")
	service.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expected)

	handler, err := search.NewSearchEventHandler(service)
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

func Test_HandleIssueTitleUpdated_ValidEvent_UpdatesIssueDocumentTitle(t *testing.T) {
	service := &mockIssueDocumentService{}
	issueID := uuid.New()
	service.On("UpdateTitle", mock.Anything, issueID, "new title").Return(nil)

	handler, err := search.NewSearchEventHandler(service)
	require.NoError(t, err)

	err = handler.HandleIssueTitleUpdated(context.Background(), trackermodel.IssueTitleUpdatedEvent{
		OccurredAt: time.Now(),
		Payload: trackermodel.Issue{
			ID:    issueID,
			Title: "new title",
		},
	})

	require.NoError(t, err)
	service.AssertExpectations(t)
}

func Test_HandleIssueTitleUpdated_ServiceError_ReturnsError(t *testing.T) {
	service := &mockIssueDocumentService{}
	expected := errors.New("update title failed")
	service.On("UpdateTitle", mock.Anything, mock.Anything, mock.Anything).Return(expected)

	handler, err := search.NewSearchEventHandler(service)
	require.NoError(t, err)

	err = handler.HandleIssueTitleUpdated(context.Background(), trackermodel.IssueTitleUpdatedEvent{
		OccurredAt: time.Now(),
		Payload: trackermodel.Issue{
			ID:    uuid.New(),
			Title: "new title",
		},
	})

	require.ErrorIs(t, err, expected)
}

func Test_HandleIssueDescriptionUpdated_ValidEvent_UpdatesIssueDocumentDescription(t *testing.T) {
	service := &mockIssueDocumentService{}
	issueID := uuid.New()
	description := "new description"
	service.On("UpdateDescription", mock.Anything, issueID, description).Return(nil)

	handler, err := search.NewSearchEventHandler(service)
	require.NoError(t, err)

	err = handler.HandleIssueDescriptionUpdated(context.Background(), trackermodel.IssueDescriptionUpdatedEvent{
		OccurredAt: time.Now(),
		Payload: trackermodel.Issue{
			ID:          issueID,
			Description: &description,
		},
	})

	require.NoError(t, err)
	service.AssertExpectations(t)
}

func Test_HandleIssueDescriptionUpdated_NilDescription_UpdatesWithEmptyDescription(t *testing.T) {
	service := &mockIssueDocumentService{}
	issueID := uuid.New()
	service.On("UpdateDescription", mock.Anything, issueID, "").Return(nil)

	handler, err := search.NewSearchEventHandler(service)
	require.NoError(t, err)

	err = handler.HandleIssueDescriptionUpdated(context.Background(), trackermodel.IssueDescriptionUpdatedEvent{
		OccurredAt: time.Now(),
		Payload: trackermodel.Issue{
			ID: issueID,
		},
	})

	require.NoError(t, err)
	service.AssertExpectations(t)
}

func Test_HandleIssueDescriptionUpdated_ServiceError_ReturnsError(t *testing.T) {
	service := &mockIssueDocumentService{}
	expected := errors.New("update description failed")
	service.On("UpdateDescription", mock.Anything, mock.Anything, mock.Anything).Return(expected)

	handler, err := search.NewSearchEventHandler(service)
	require.NoError(t, err)

	err = handler.HandleIssueDescriptionUpdated(context.Background(), trackermodel.IssueDescriptionUpdatedEvent{
		OccurredAt: time.Now(),
		Payload: trackermodel.Issue{
			ID: uuid.New(),
		},
	})

	require.ErrorIs(t, err, expected)
}
