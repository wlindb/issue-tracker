package tracker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	"github.com/wlindb/issue-tracker/internal/domain/tracker/comment"
	"github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
	"github.com/wlindb/issue-tracker/internal/domain/tracker/project"
	key "github.com/wlindb/issue-tracker/internal/pkg/context"
	embeddednats "github.com/wlindb/issue-tracker/internal/pkg/nats"
)

type eventPublisher struct {
	natsEventPublisher *embeddednats.NATSEventPublisher
}

func NewEventPublisher(connection *nats.Conn) error {
	publisher := eventPublisher{
		natsEventPublisher: embeddednats.NewNATSEventPublisher(connection),
	}

	if err := issue.Created.AddPublisher(publisher.IssueCreatedPublisher); err != nil {
		return fmt.Errorf("issue created event publisher: %w", err)
	}

	if err := issue.StatusUpdated.AddPublisher(publisher.IssueStatusUpdatedPublisher); err != nil {
		return fmt.Errorf("issue status updated event publisher: %w", err)
	}

	if err := issue.TitleUpdated.AddPublisher(publisher.IssueTitleUpdatedPublisher); err != nil {
		return fmt.Errorf("issue title updated event publisher: %w", err)
	}

	if err := issue.PriorityUpdated.AddPublisher(publisher.IssuePriorityUpdatedPublisher); err != nil {
		return fmt.Errorf("issue priority updated event publisher: %w", err)
	}

	if err := issue.AssigneeUpdated.AddPublisher(publisher.IssueAssigneeUpdatedPublisher); err != nil {
		return fmt.Errorf("issue assignee updated event publisher: %w", err)
	}

	if err := issue.DescriptionUpdated.AddPublisher(publisher.IssueDescriptionUpdatedPublisher); err != nil {
		return fmt.Errorf("issue description updated event publisher: %w", err)
	}

	if err := issue.LabelAdded.AddPublisher(publisher.IssueLabelAddedPublisher); err != nil {
		return fmt.Errorf("issue label added event publisher: %w", err)
	}

	if err := comment.Created.AddPublisher(publisher.CommentCreatedPublisher); err != nil {
		return fmt.Errorf("comment created event publisher: %w", err)
	}

	if err := project.Created.AddPublisher(publisher.ProjectCreatedPublisher); err != nil {
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

func (publisher eventPublisher) IssueCreatedPublisher(ctx context.Context, event issue.IssueCreatedEvent) error {
	payload, err := json.Marshal(ToIssueCreatedEventDTO(event))
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	resolver := IssueCreatedSubjectResolver{}
	subject, err := resolver.Resolve(ctx, event)
	if err != nil {
		return fmt.Errorf("resolve: %w", err)
	}

	if err := publisher.natsEventPublisher.Publish(ctx, subject, payload); err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	return nil
}

type IssueStatusUpdatedSubjectResolver struct{}

func (IssueStatusUpdatedSubjectResolver) Resolve(ctx context.Context, _ issue.IssueStatusUpdatedEvent) (string, error) {
	workspaceID, err := workspaceID(ctx)
	if err != nil {
		return "", fmt.Errorf("issue status updated subject resolver: %w", err)
	}
	return embeddednats.IssueStatusUpdatedSubject.Subject(workspaceID), nil
}

func (publisher eventPublisher) IssueStatusUpdatedPublisher(ctx context.Context, event issue.IssueStatusUpdatedEvent) error {
	payload, err := json.Marshal(ToIssueStatusUpdatedEvent(event))
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	resolver := IssueStatusUpdatedSubjectResolver{}
	subject, err := resolver.Resolve(ctx, event)
	if err != nil {
		return fmt.Errorf("resolve: %w", err)
	}

	if err := publisher.natsEventPublisher.Publish(ctx, subject, payload); err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	return nil
}

type IssueTitleUpdatedSubjectResolver struct{}

func (IssueTitleUpdatedSubjectResolver) Resolve(ctx context.Context, _ issue.IssueTitleUpdatedEvent) (string, error) {
	workspaceID, err := workspaceID(ctx)
	if err != nil {
		return "", fmt.Errorf("issue title updated subject resolver: %w", err)
	}
	return embeddednats.IssueTitleUpdatedSubject.Subject(workspaceID), nil
}

func (publisher eventPublisher) IssueTitleUpdatedPublisher(ctx context.Context, event issue.IssueTitleUpdatedEvent) error {
	payload, err := json.Marshal(ToIssueTitleUpdatedEvent(event))
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	resolver := IssueTitleUpdatedSubjectResolver{}
	subject, err := resolver.Resolve(ctx, event)
	if err != nil {
		return fmt.Errorf("resolve: %w", err)
	}

	if err := publisher.natsEventPublisher.Publish(ctx, subject, payload); err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	return nil
}

type IssuePriorityUpdatedSubjectResolver struct{}

func (IssuePriorityUpdatedSubjectResolver) Resolve(ctx context.Context, _ issue.IssuePriorityUpdatedEvent) (string, error) {
	workspaceID, err := workspaceID(ctx)
	if err != nil {
		return "", fmt.Errorf("issue priority updated subject resolver: %w", err)
	}
	return embeddednats.IssuePriorityUpdatedSubject.Subject(workspaceID), nil
}

func (publisher eventPublisher) IssuePriorityUpdatedPublisher(ctx context.Context, event issue.IssuePriorityUpdatedEvent) error {
	payload, err := json.Marshal(ToIssuePriorityUpdatedEvent(event))
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	resolver := IssuePriorityUpdatedSubjectResolver{}
	subject, err := resolver.Resolve(ctx, event)
	if err != nil {
		return fmt.Errorf("resolve: %w", err)
	}

	if err := publisher.natsEventPublisher.Publish(ctx, subject, payload); err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	return nil
}

type IssueAssigneeUpdatedSubjectResolver struct{}

func (IssueAssigneeUpdatedSubjectResolver) Resolve(ctx context.Context, _ issue.IssueAssigneeUpdatedEvent) (string, error) {
	workspaceID, err := workspaceID(ctx)
	if err != nil {
		return "", fmt.Errorf("issue assignee updated subject resolver: %w", err)
	}
	return embeddednats.IssueAssigneeUpdatedSubject.Subject(workspaceID), nil
}

func (publisher eventPublisher) IssueAssigneeUpdatedPublisher(ctx context.Context, event issue.IssueAssigneeUpdatedEvent) error {
	payload, err := json.Marshal(ToIssueAssigneeUpdatedEvent(event))
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	resolver := IssueAssigneeUpdatedSubjectResolver{}
	subject, err := resolver.Resolve(ctx, event)
	if err != nil {
		return fmt.Errorf("resolve: %w", err)
	}

	if err := publisher.natsEventPublisher.Publish(ctx, subject, payload); err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	return nil
}

type IssueDescriptionUpdatedSubjectResolver struct{}

func (IssueDescriptionUpdatedSubjectResolver) Resolve(ctx context.Context, _ issue.IssueDescriptionUpdatedEvent) (string, error) {
	workspaceID, err := workspaceID(ctx)
	if err != nil {
		return "", fmt.Errorf("issue description updated subject resolver: %w", err)
	}
	return embeddednats.IssueDescriptionUpdatedSubject.Subject(workspaceID), nil
}

func (publisher eventPublisher) IssueDescriptionUpdatedPublisher(ctx context.Context, event issue.IssueDescriptionUpdatedEvent) error {
	payload, err := json.Marshal(ToIssueDescriptionUpdatedEvent(event))
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	resolver := IssueDescriptionUpdatedSubjectResolver{}
	subject, err := resolver.Resolve(ctx, event)
	if err != nil {
		return fmt.Errorf("resolve: %w", err)
	}

	if err := publisher.natsEventPublisher.Publish(ctx, subject, payload); err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	return nil
}

type IssueLabelAddedSubjectResolver struct{}

func (IssueLabelAddedSubjectResolver) Resolve(ctx context.Context, _ issue.IssueLabelAddedEvent) (string, error) {
	workspaceID, err := workspaceID(ctx)
	if err != nil {
		return "", fmt.Errorf("issue label added subject resolver: %w", err)
	}
	return embeddednats.IssueLabelAddedSubject.Subject(workspaceID), nil
}

func (publisher eventPublisher) IssueLabelAddedPublisher(ctx context.Context, event issue.IssueLabelAddedEvent) error {
	payload, err := json.Marshal(ToIssueLabelAddedEvent(event))
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	resolver := IssueLabelAddedSubjectResolver{}
	subject, err := resolver.Resolve(ctx, event)
	if err != nil {
		return fmt.Errorf("resolve: %w", err)
	}

	if err := publisher.natsEventPublisher.Publish(ctx, subject, payload); err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	return nil
}

type CommentCreatedSubjectResolver struct{}

func (CommentCreatedSubjectResolver) Resolve(ctx context.Context, event comment.CommentCreatedEvent) (string, error) {
	workspaceID, err := workspaceID(ctx)
	if err != nil {
		return "", fmt.Errorf("comment publisher subject resolver: %w", err)
	}

	return embeddednats.CommentCreatedSubject.Subject(workspaceID, event.Payload.IssueID), nil
}

func (publisher eventPublisher) CommentCreatedPublisher(ctx context.Context, event comment.CommentCreatedEvent) error {
	payload, err := json.Marshal(ToCommentCreatedEvent(event))
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	resolver := CommentCreatedSubjectResolver{}
	subject, err := resolver.Resolve(ctx, event)
	if err != nil {
		return fmt.Errorf("resolve: %w", err)
	}

	if err := publisher.natsEventPublisher.Publish(ctx, subject, payload); err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	return nil
}

type ProjectCreatedSubjectResolver struct{}

func (ProjectCreatedSubjectResolver) Resolve(ctx context.Context, _ project.ProjectCreatedEvent) (string, error) {
	workspaceID, err := workspaceID(ctx)
	if err != nil {
		return "", fmt.Errorf("issue created subject resolver: %w", err)
	}

	return embeddednats.ProjectCreatedSubject.Subject(workspaceID), nil
}

func (publisher eventPublisher) ProjectCreatedPublisher(ctx context.Context, event project.ProjectCreatedEvent) error {
	payload, err := json.Marshal(ToProjectCreatedEvent(event))
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	resolver := ProjectCreatedSubjectResolver{}
	subject, err := resolver.Resolve(ctx, event)
	if err != nil {
		return fmt.Errorf("resolve: %w", err)
	}

	if err := publisher.natsEventPublisher.Publish(ctx, subject, payload); err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	return nil
}

func workspaceID(ctx context.Context) (uuid.UUID, error) {
	workspaceID, ok := ctx.Value(key.WorkspaceID).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("workspace ID missing from context")
	}
	return workspaceID, nil
}
