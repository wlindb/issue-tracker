package tracker

import (
	"fmt"

	"github.com/nats-io/nats.go"

	"github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
	"github.com/wlindb/issue-tracker/internal/domain/tracker/project"
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

	projectPublisher := embeddednats.NewNATSEventPublisher[project.ProjectCreatedEvent](
		connection,
		embeddednats.ProjectCreatedSubject,
	)
	if err := project.Created.AddPublisher(projectPublisher.Publisher); err != nil {
		return fmt.Errorf("project created event publisher: %w", err)
	}

	return nil
}
