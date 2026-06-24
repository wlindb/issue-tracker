package issue

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

// IssueService implements the domain logic for managing issues.
type IssueService struct {
	unitOfWork UnitOfWork
	repository IssueRepository
}

type Repositories struct {
	Issues IssueRepository
}

type UnitOfWork interface {
	RunInTx(ctx context.Context, fn func(Repositories) error) error
}

// NewIssueService creates an IssueService wired to the given repository and event publisher.
func NewIssueService(
	unitOfWork UnitOfWork,
	repository IssueRepository,
) *IssueService {
	return &IssueService{unitOfWork: unitOfWork, repository: repository}
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
	var issue Issue
	if err := s.unitOfWork.RunInTx(ctx, func(tx Repositories) error {
		var err error
		issue, err = tx.Issues.CreateIssue(ctx, command.ToIssue(uuid.New(), command.Slugify, []Label{}))
		if err != nil {
			return fmt.Errorf("create issue: %w", err)
		}

		if err := issue.EmitCreated(ctx); err != nil {
			slog.Error("publish issue created event", "error", err)
		}

		return nil
	}); err != nil {
		return Issue{}, fmt.Errorf("run in tx: %w", err)
	}

	return issue, nil
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

// UpdateIssueTitle updates the title of an issue.
func (s *IssueService) UpdateIssueTitle(ctx context.Context, issueID uuid.UUID, title string) (Issue, error) {
	current, err := s.repository.GetIssue(ctx, issueID)
	if err != nil {
		return Issue{}, fmt.Errorf("update issue title: %w", err)
	}
	updated, err := current.UpdateTitle(title)
	if err != nil {
		return Issue{}, fmt.Errorf("update issue title: %w", err)
	}
	result, err := s.repository.Update(ctx, updated)
	if err != nil {
		return Issue{}, fmt.Errorf("update issue title: %w", err)
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
	if err := result.EmitStatusUpdated(ctx); err != nil {
		slog.Error("publish issue status updated event", "error", err)
	}
	return result, nil
}
