package tracker

import (
	"context"
	"encoding/json"

	natsgo "github.com/nats-io/nats.go"

	"github.com/wlindb/issue-tracker/internal/pkg/domain/event"
)

const issueCreatedSubject = "issues.created"

// NATSEventPublisher implements issuedomain.EventPublisher using a NATS connection.
type NATSEventPublisher[T any] struct {
	connection *natsgo.Conn
	Publisher  event.Publisher[T]
}

// NewNATSEventPublisher returns a NATSEventPublisher that publishes to the given connection.
func NewNATSEventPublisher[T any](connection *natsgo.Conn) *NATSEventPublisher[T] {
	return &NATSEventPublisher[T]{connection: connection, Publisher: publish[T](connection)}
}

// PublishIssueCreated marshals the event as JSON and publishes it to the "issues.created" subject.
// func (p *NATSEventPublisher) PublishIssueCreated(event issuedomain.IssueCreatedEvent) error {
// 	payload, err := json.Marshal(event)
// 	if err != nil {
// 		return fmt.Errorf("marshal issue created event: %w", err)
// 	}
// 	if err := p.connection.Publish(issueCreatedSubject, payload); err != nil {
// 		return fmt.Errorf("publish issue created event: %w", err)
// 	}
// 	return nil
// }

func publish[T any](connection *natsgo.Conn) event.Publisher[T] {
	return func(_ context.Context, event T) {
		payload, err := json.Marshal(event)
		if err != nil {
			return
		}
		if err := connection.Publish(issueCreatedSubject, payload); err != nil {
			return
		}
	}
}
