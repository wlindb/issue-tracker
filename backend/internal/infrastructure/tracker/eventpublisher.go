package tracker

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	"github.com/wlindb/issue-tracker/internal/domain/tracker/comment"
	"github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
	key "github.com/wlindb/issue-tracker/internal/pkg/context"
	"github.com/wlindb/issue-tracker/internal/domain/tracker/project"
	embeddednats "github.com/wlindb/issue-tracker/internal/pkg/nats"
)

func NewEventPublisher(connection *nats.Conn) error {
	issuePublisher := embeddednats.NewNATSEventPublisher(
		connection,
		IssueCreatedSubjectResolver{},
	)
	if err := issue.Created.AddPublisher(issuePublisher.Publisher); err != nil {
		return fmt.Errorf("issue created event publisher: %w", err)
	}

	commentPublisher := embeddednats.NewNATSEventPublisher(
		connection,
		CommentCreatedSubjectResolver{},
	)
	if err := comment.Created.AddPublisher(commentPublisher.Publisher); err != nil {
		return fmt.Errorf("comment created event publisher: %w", err)
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

type IssueCreatedSubjectResolver struct{}

func (IssueCreatedSubjectResolver) Resolve(ctx context.Context, _ issue.IssueCreatedEvent) (string, error) {
	workspaceID, err := workspaceID(ctx)
	if err != nil {
		return "", fmt.Errorf("issue created subject resolver: %w", err)
	}

	return embeddednats.IssueCreatedSubject.Subject(workspaceID), nil
}

type CommentCreatedSubjectResolver struct{}

func (CommentCreatedSubjectResolver) Resolve(ctx context.Context, event comment.CommentCreatedEvent) (string, error) {
	workspaceID, err := workspaceID(ctx)
	if err != nil {
		return "", fmt.Errorf("comment publisher subject resolver: %w", err)
	}

	return embeddednats.CommentCreatedSubject.Subject(workspaceID, event.Payload.IssueID), nil
}

func workspaceID(ctx context.Context) (uuid.UUID, error) {
	workspaceID, ok := ctx.Value(key.WorkspaceID).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("workspace ID missing from context")
	}
	return workspaceID, nil
}
