package tracker

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	"github.com/wlindb/issue-tracker/internal/domain/tracker/comment"
	"github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
	key "github.com/wlindb/issue-tracker/internal/pkg/context"
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

	commentPublisher := embeddednats.NewNATSIssueEventPublisher(
		connection,
		func(ctx context.Context, event comment.CommentCreatedEvent) (string, error) {
			workspaceID, ok := ctx.Value(key.WorkspaceID).(uuid.UUID)
			if !ok {
				return "", fmt.Errorf("comment created subject: workspace ID missing from context")
			}

			return embeddednats.CommentCreatedSubject.Subject(workspaceID, event.Payload.IssueID), nil
		},
	)
	if err := comment.Created.AddPublisher(commentPublisher.Publisher); err != nil {
		return fmt.Errorf("comment created event publisher: %w", err)
	}

	return nil
}
