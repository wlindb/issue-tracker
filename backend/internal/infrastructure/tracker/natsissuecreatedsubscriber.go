package tracker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	natsgo "github.com/nats-io/nats.go"

	applicationtracker "github.com/wlindb/issue-tracker/internal/application/tracker"
	issuedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
	"github.com/wlindb/issue-tracker/internal/pkg/domain/event"
)

// NATSIssueCreatedSubscriber wires the NATS transport to the application-layer handler.
type NATSIssueCreatedSubscriber struct {
	connection *natsgo.Conn
	handler    *applicationtracker.IssueCreatedHandler
}

// NewNATSIssueCreatedSubscriber returns a NATSIssueCreatedSubscriber.
func NewNATSIssueCreatedSubscriber(connection *natsgo.Conn, handler *applicationtracker.IssueCreatedHandler) *NATSIssueCreatedSubscriber {
	return &NATSIssueCreatedSubscriber{connection: connection, handler: handler}
}

// Subscribe registers the handler on the "issues.created" subject.
// The returned subscription must be drained or unsubscribed by the caller on shutdown.
func (s *NATSIssueCreatedSubscriber) Subscribe() (*natsgo.Subscription, error) {
	subscription, err := s.connection.Subscribe(issueCreatedSubject, func(message *natsgo.Msg) {
		var event issuedomain.IssueCreatedEvent
		if err := json.Unmarshal(message.Data, &event); err != nil {
			slog.Error("unmarshal issue created event", "error", err)
			return
		}
		s.handler.Handle(event)
	})
	if err != nil {
		return nil, fmt.Errorf("subscribe issue created: %w", err)
	}
	return subscription, nil
}

type NATSEventSubscriber[T any] struct {
	connection *natsgo.Conn
	Subscriber event.Subscriber[T]
}

func NewNATSEventSubscriber[T any](connection *natsgo.Conn) NATSEventSubscriber[T] {
	return NATSEventSubscriber[T]{connection: connection}
}

func (s NATSEventSubscriber[T]) Subscribe(handler event.Subscriber[T]) error {
	_, err := s.connection.Subscribe(issueCreatedSubject, func(message *natsgo.Msg) {
		var event T
		if err := json.Unmarshal(message.Data, &event); err != nil {
			slog.Error("unmarshal issue created event", "error", err)
			return
		}
		// TODO: Fix context
		ctx := context.Background()
		handler(ctx, event)
	})
	if err != nil {
		return fmt.Errorf("subscribe subject %s: %w", issueCreatedSubject, err)
	}

	return nil
}
