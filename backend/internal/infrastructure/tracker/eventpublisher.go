package tracker

import (
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
	embeddednats "github.com/wlindb/issue-tracker/internal/pkg/nats"
)

func NewEventPublisher(connection *nats.Conn) error {
	publisher := embeddednats.NewNATSEventPublisher[issue.IssueCreatedEvent](connection, embeddednats.IssueCreatedSubject)
	if err := issue.Created.AddPublisher(publisher.Publisher); err != nil {
		return fmt.Errorf("issue created event publisher: %w", err)
	}
	return nil
}
