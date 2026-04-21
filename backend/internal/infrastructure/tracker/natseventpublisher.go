package tracker

import (
	"context"
	"encoding/json"
	"fmt"

	natsgo "github.com/nats-io/nats.go"

	"github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
	issuedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
	"github.com/wlindb/issue-tracker/internal/pkg/domain/event"
)

const issueCreatedSubject = "issues.created"

// NATSEventPublisher implements issuedomain.EventPublisher using a NATS connection.
type NATSEventPublisher struct {
	connection *natsgo.Conn
	Publisher  event.Subscriber[issue.IssueCreatedEvent]
}

// NewNATSEventPublisher returns a NATSEventPublisher that publishes to the given connection.
func NewNATSEventPublisher(connection *natsgo.Conn) *NATSEventPublisher {
	return &NATSEventPublisher{connection: connection, Publisher: publish(connection)}
}

// PublishIssueCreated marshals the event as JSON and publishes it to the "issues.created" subject.
func (p *NATSEventPublisher) PublishIssueCreated(event issuedomain.IssueCreatedEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal issue created event: %w", err)
	}
	if err := p.connection.Publish(issueCreatedSubject, payload); err != nil {
		return fmt.Errorf("publish issue created event: %w", err)
	}
	return nil
}

func publish(connection *natsgo.Conn) event.Subscriber[issue.IssueCreatedEvent] {
	return func(_ context.Context, event issuedomain.IssueCreatedEvent) {
		payload, err := json.Marshal(event)
		if err != nil {
			return
		}
		if err := connection.Publish(issueCreatedSubject, payload); err != nil {
			return
		}
	}
}
