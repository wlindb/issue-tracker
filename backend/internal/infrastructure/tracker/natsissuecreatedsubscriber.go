package tracker

import (
	"encoding/json"
	"fmt"
	"log/slog"

	natsgo "github.com/nats-io/nats.go"

	applicationtracker "github.com/wlindb/issue-tracker/internal/application/tracker"
	issuedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
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
