package issue

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	label "github.com/wlindb/issue-tracker/internal/domain/tracker/label"
)

// IssueService implements the domain logic for managing issues.
type IssueService struct {
	unitOfWork  UnitOfWork
	repository  IssueRepository
	labelLister LabelLister
}

type Repositories struct {
	Issues IssueRepository
}

type UnitOfWork interface {
	RunInTx(ctx context.Context, fn func(Repositories) error) error
}

type LabelLister interface {
	ListByIDs(ctx context.Context, ids []uuid.UUID) ([]label.Label, error)
}

// NewIssueService creates an IssueService wired to the given repository and event publisher.
func NewIssueService(
	unitOfWork UnitOfWork,
	repository IssueRepository,
	labelLister LabelLister,
) *IssueService {
	return &IssueService{unitOfWork: unitOfWork, repository: repository, labelLister: labelLister}
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
	labels, err := s.labelLister.ListByIDs(ctx, command.LabelIDs)
	if err != nil {
		return Issue{}, fmt.Errorf("create issue: %w", err)
	}

	var issue Issue
	if err := s.unitOfWork.RunInTx(ctx, func(tx Repositories) error {
		createdIssue, err := tx.Issues.CreateIssue(ctx, command.ToIssue(uuid.New(), command.Slugify, labels))
		if err != nil {
			return fmt.Errorf("create issue: %w", err)
		}
		issue = createdIssue

		if err := issue.EmitCreated(ctx); err != nil {
			return fmt.Errorf("emit created: %w", err)
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

	var issue Issue
	if err := s.unitOfWork.RunInTx(ctx, func(tx Repositories) error {
		var err error
		issue, err = tx.Issues.Update(ctx, updated)
		if err != nil {
			return fmt.Errorf("update issue assignee: %w", err)
		}
		if err := issue.EmitAssigneeUpdated(ctx); err != nil {
			return fmt.Errorf("emit assignee updated: %w", err)
		}

		return nil
	}); err != nil {
		return Issue{}, fmt.Errorf("run in tx: %w", err)
	}
	return issue, nil
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

	var result Issue
	if err := s.unitOfWork.RunInTx(ctx, func(tx Repositories) error {
		var err error
		result, err = tx.Issues.Update(ctx, updated)
		if err != nil {
			return fmt.Errorf("update issue description: %w", err)
		}
		if err := result.EmitDescriptionUpdated(ctx); err != nil {
			return fmt.Errorf("emit description updated: %w", err)
		}

		return nil
	}); err != nil {
		return Issue{}, fmt.Errorf("run in tx: %w", err)
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
	var result Issue
	if err := s.unitOfWork.RunInTx(ctx, func(tx Repositories) error {
		var err error
		result, err = tx.Issues.Update(ctx, updated)
		if err != nil {
			return fmt.Errorf("update issue title: %w", err)
		}
		if err := result.EmitTitleUpdated(ctx); err != nil {
			return fmt.Errorf("emit title updated: %w", err)
		}
		return nil
	}); err != nil {
		return Issue{}, fmt.Errorf("run in tx: %w", err)
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

	var result Issue
	if err = s.unitOfWork.RunInTx(ctx, func(tx Repositories) error {
		result, err = tx.Issues.Update(ctx, updated)
		if err != nil {
			return fmt.Errorf("update issue priority: %w", err)
		}
		if err = result.EmitPriorityUpdated(ctx); err != nil {
			return fmt.Errorf("emit priority updated: %w", err)
		}
		return nil
	}); err != nil {
		return Issue{}, fmt.Errorf("run in tx: %w", err)
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

// AddLabel attaches a label to an issue. Adding a label the issue already has is a no-op.
func (s *IssueService) AddLabel(ctx context.Context, issueID uuid.UUID, label label.Label) (Issue, error) {
	current, err := s.repository.GetIssue(ctx, issueID)
	if err != nil {
		return Issue{}, fmt.Errorf("add label: %w", err)
	}

	if current.HasLabel(label) {
		return current, nil
	}

	updated := current.AddLabel(label)

	var result Issue
	if err := s.unitOfWork.RunInTx(ctx, func(tx Repositories) error {
		if err := tx.Issues.AddLabel(ctx, issueID, label); err != nil {
			return fmt.Errorf("add label: %w", err)
		}
		result = updated
		if err := result.EmitLabelAdded(ctx); err != nil {
			return fmt.Errorf("emit label added: %w", err)
		}
		return nil
	}); err != nil {
		return Issue{}, fmt.Errorf("run in tx: %w", err)
	}

	return result, nil
}
