//go:build !integration

package search_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
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

func Test_NewSearchEventHandler_WithIssueCreated_SubscribesHandleIssueCreated(t *testing.T) {
	subscriber := &fakeSubscriber[trackermodel.IssueCreatedEvent]{}

	_, err := search.NewSearchEventHandler(search.WithIssueCreated(subscriber))
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

	_, err := search.NewSearchEventHandler(search.WithIssueTitleUpdated(subscriber))
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

	_, err := search.NewSearchEventHandler(search.WithIssueDescriptionUpdated(subscriber))
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

	_, err := search.NewSearchEventHandler(search.WithIssueCreated(subscriber))
	require.ErrorIs(t, err, expected)
}
