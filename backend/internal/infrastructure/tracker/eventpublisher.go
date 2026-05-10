package tracker

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	"github.com/wlindb/issue-tracker/internal/domain/tracker/comment"
	"github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
	embeddednats "github.com/wlindb/issue-tracker/internal/pkg/nats"
)

func NewEventPublisher(connection *nats.Conn) error {
	issuePublisher := embeddednats.NewNATSEventPublisher[issue.IssueCreatedEvent](
		connection,
		embeddednats.IssueCreatedSubject,
	)
	if err := issue.Created.AddPublisher(issuePublisher.Publisher); err != nil {
		return fmt.Errorf("issue created event publisher: %w", err)
	}

	commentPublisher := embeddednats.NewNATSIssueEventPublisher[comment.CommentCreatedEvent](
		connection,
		embeddednats.CommentCreatedSubject,
		func(evt comment.CommentCreatedEvent) uuid.UUID { return evt.Payload.IssueID },
	)
	if err := comment.Created.AddPublisher(commentPublisher.Publisher); err != nil {
		return fmt.Errorf("comment created event publisher: %w", err)
	}

	return nil
}
