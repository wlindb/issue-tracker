package tracker

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	"github.com/wlindb/issue-tracker/internal/domain/tracker/comment"
	"github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
	"github.com/wlindb/issue-tracker/internal/domain/tracker/project"
	key "github.com/wlindb/issue-tracker/internal/pkg/context"
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

	statusUpdatedPublisher := embeddednats.NewNATSEventPublisher(
		connection,
		IssueStatusUpdatedSubjectResolver{},
	)
	if err := issue.StatusUpdated.AddPublisher(statusUpdatedPublisher.Publisher); err != nil {
		return fmt.Errorf("issue status updated event publisher: %w", err)
	}

	titleUpdatedPublisher := embeddednats.NewNATSEventPublisher(
		connection,
		IssueTitleUpdatedSubjectResolver{},
	)
	if err := issue.TitleUpdated.AddPublisher(titleUpdatedPublisher.Publisher); err != nil {
		return fmt.Errorf("issue title updated event publisher: %w", err)
	}

	priorityUpdatedPublisher := embeddednats.NewNATSEventPublisher(
		connection,
		IssuePriorityUpdatedSubjectResolver{},
	)
	if err := issue.PriorityUpdated.AddPublisher(priorityUpdatedPublisher.Publisher); err != nil {
		return fmt.Errorf("issue priority updated event publisher: %w", err)
	}

	commentPublisher := embeddednats.NewNATSEventPublisher(
		connection,
		CommentCreatedSubjectResolver{},
	)
	if err := comment.Created.AddPublisher(commentPublisher.Publisher); err != nil {
		return fmt.Errorf("comment created event publisher: %w", err)
	}

	projectPublisher := embeddednats.NewNATSEventPublisher(
		connection,
		ProjectCreatedSubjectResolver{},
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

type IssueStatusUpdatedSubjectResolver struct{}

func (IssueStatusUpdatedSubjectResolver) Resolve(ctx context.Context, _ issue.IssueStatusUpdatedEvent) (string, error) {
	workspaceID, err := workspaceID(ctx)
	if err != nil {
		return "", fmt.Errorf("issue status updated subject resolver: %w", err)
	}
	return embeddednats.IssueStatusUpdatedSubject.Subject(workspaceID), nil
}

type IssueTitleUpdatedSubjectResolver struct{}

func (IssueTitleUpdatedSubjectResolver) Resolve(ctx context.Context, _ issue.IssueTitleUpdatedEvent) (string, error) {
	workspaceID, err := workspaceID(ctx)
	if err != nil {
		return "", fmt.Errorf("issue title updated subject resolver: %w", err)
	}
	return embeddednats.IssueTitleUpdatedSubject.Subject(workspaceID), nil
}

type IssuePriorityUpdatedSubjectResolver struct{}

func (IssuePriorityUpdatedSubjectResolver) Resolve(ctx context.Context, _ issue.IssuePriorityUpdatedEvent) (string, error) {
	workspaceID, err := workspaceID(ctx)
	if err != nil {
		return "", fmt.Errorf("issue priority updated subject resolver: %w", err)
	}
	return embeddednats.IssuePriorityUpdatedSubject.Subject(workspaceID), nil
}

type CommentCreatedSubjectResolver struct{}

func (CommentCreatedSubjectResolver) Resolve(ctx context.Context, event comment.CommentCreatedEvent) (string, error) {
	workspaceID, err := workspaceID(ctx)
	if err != nil {
		return "", fmt.Errorf("comment publisher subject resolver: %w", err)
	}

	return embeddednats.CommentCreatedSubject.Subject(workspaceID, event.Payload.IssueID), nil
}

type ProjectCreatedSubjectResolver struct{}

func (ProjectCreatedSubjectResolver) Resolve(ctx context.Context, _ project.ProjectCreatedEvent) (string, error) {
	workspaceID, err := workspaceID(ctx)
	if err != nil {
		return "", fmt.Errorf("issue created subject resolver: %w", err)
	}

	return embeddednats.ProjectCreatedSubject.Subject(workspaceID), nil
}

func workspaceID(ctx context.Context) (uuid.UUID, error) {
	workspaceID, ok := ctx.Value(key.WorkspaceID).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("workspace ID missing from context")
	}
	return workspaceID, nil
}
