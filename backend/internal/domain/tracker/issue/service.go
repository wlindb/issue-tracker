package issue

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// IssueService implements the domain logic for managing issues.
type IssueService struct {
	repository IssueRepository
	publisher  EventPublisher
}

// NewIssueService creates an IssueService wired to the given repository and event publisher.
func NewIssueService(repository IssueRepository, publisher EventPublisher) *IssueService {
	return &IssueService{repository: repository, publisher: publisher}
}

// ListIssues returns a paginated list of issues for the given project.
func (s *IssueService) ListIssues(ctx context.Context, projectID uuid.UUID, query ListIssueQuery) (IssuePage, error) {
	page, err := s.repository.ListIssues(ctx, projectID, query)
	if err != nil {
		return IssuePage{}, fmt.Errorf("list issues: %w", err)
	}
	return page, nil
}

// GetIssue retrieves a single issue by its ID.
func (s *IssueService) GetIssue(ctx context.Context, issueID uuid.UUID) (Issue, error) {
	issue, err := s.repository.GetIssue(ctx, issueID)
	if err != nil {
		return Issue{}, fmt.Errorf("get issue: %w", err)
	}
	return issue, nil
}

// CreateIssue creates a new issue from the given command.
func (s *IssueService) CreateIssue(ctx context.Context, command CreateIssueCommand) (Issue, error) {
	issue := command.ToIssue(uuid.New(), command.Slugify)
	result, err := s.repository.CreateIssue(ctx, issue)
	if err != nil {
		return Issue{}, fmt.Errorf("create issue: %w", err)
	}
	if publishErr := s.publisher.PublishIssueCreated(IssueCreatedEvent{
		IssueID:    result.ID,
		ProjectID:  result.ProjectID,
		ReporterID: result.ReporterID,
		Title:      result.Title,
		Status:     result.Status,
		Priority:   result.Priority,
		OccurredAt: time.Now().UTC(),
	}); publishErr != nil {
		slog.Error("publish issue created event", "error", publishErr)
	}
	return result, nil
}

// UpdateIssueAssignee updates the assignee of an issue.
func (s *IssueService) UpdateIssueAssignee(ctx context.Context, issueID uuid.UUID, assigneeID *uuid.UUID) (Issue, error) {
	current, err := s.repository.GetIssue(ctx, issueID)
	if err != nil {
		return Issue{}, fmt.Errorf("update issue assignee: %w", err)
	}
	var updated Issue
	if assigneeID != nil {
		updated, err = current.UpdateAssignee(*assigneeID)
		if err != nil {
			return Issue{}, fmt.Errorf("update issue assignee: %w", err)
		}
	} else {
		updated = current.Unassign()
	}
	result, err := s.repository.Update(ctx, updated)
	if err != nil {
		return Issue{}, fmt.Errorf("update issue assignee: %w", err)
	}
	return result, nil
}

// UpdateIssueDescription updates the description of an issue.
func (s *IssueService) UpdateIssueDescription(ctx context.Context, issueID uuid.UUID, description *string) (Issue, error) {
	current, err := s.repository.GetIssue(ctx, issueID)
	if err != nil {
		return Issue{}, fmt.Errorf("update issue description: %w", err)
	}
	updated, err := current.UpdateDescription(description)
	if err != nil {
		return Issue{}, fmt.Errorf("update issue description: %w", err)
	}
	result, err := s.repository.Update(ctx, updated)
	if err != nil {
		return Issue{}, fmt.Errorf("update issue description: %w", err)
	}
	return result, nil
}

// UpdateIssuePriority updates the priority of an issue.
func (s *IssueService) UpdateIssuePriority(ctx context.Context, issueID uuid.UUID, priority Priority) (Issue, error) {
	current, err := s.repository.GetIssue(ctx, issueID)
	if err != nil {
		return Issue{}, fmt.Errorf("update issue priority: %w", err)
	}
	updated, err := current.UpdatePriority(priority)
	if err != nil {
		return Issue{}, fmt.Errorf("update issue priority: %w", err)
	}
	result, err := s.repository.Update(ctx, updated)
	if err != nil {
		return Issue{}, fmt.Errorf("update issue priority: %w", err)
	}
	return result, nil
}

// UpdateIssueStatus updates the status of the issue with the given ID.
func (s *IssueService) UpdateIssueStatus(ctx context.Context, issueID uuid.UUID, status Status) (Issue, error) {
	current, err := s.repository.GetIssue(ctx, issueID)
	if err != nil {
		return Issue{}, fmt.Errorf("update issue status: %w", err)
	}
	updated, err := current.UpdateStatus(status)
	if err != nil {
		return Issue{}, fmt.Errorf("update issue status: %w", err)
	}
	result, err := s.repository.Update(ctx, updated)
	if err != nil {
		return Issue{}, fmt.Errorf("update issue status: %w", err)
	}
	return result, nil
}
