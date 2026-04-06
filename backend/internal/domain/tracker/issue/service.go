package issue

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// IssueService implements the domain logic for managing issues.
type IssueService struct {
	repository IssueRepository
}

// NewIssueService creates an IssueService wired to the given repository.
func NewIssueService(repository IssueRepository) *IssueService {
	return &IssueService{repository: repository}
}

// ListIssues returns a paginated list of issues for the given project.
func (s *IssueService) ListIssues(ctx context.Context, projectID uuid.UUID, query ListIssueQuery) (IssuePage, error) {
	page, err := s.repository.ListIssues(ctx, projectID, query)
	if err != nil {
		return IssuePage{}, fmt.Errorf("list issues: %w", err)
	}
	return page, nil
}

// UpdateIssueAssignee updates the assignee of an issue.
func (s *IssueService) UpdateIssueAssignee(_ context.Context, _ uuid.UUID, _ *uuid.UUID) (Issue, error) {
	return Issue{}, errors.New("not implemented")
}

// CreateIssue creates a new issue from the given command.
func (s *IssueService) CreateIssue(ctx context.Context, command CreateIssueCommand) (Issue, error) {
	issue := command.ToIssue(uuid.New(), command.Slugify)
	result, err := s.repository.CreateIssue(ctx, issue)
	if err != nil {
		return Issue{}, fmt.Errorf("create issue: %w", err)
	}
	return result, nil
}

// UpdateIssueDescription is not yet implemented.
func (s *IssueService) UpdateIssueDescription(_ context.Context, _ uuid.UUID, _ *string) (Issue, error) {
	return Issue{}, errors.New("not implemented")
}

func (s *IssueService) UpdateIssuePriority(_ context.Context, _ uuid.UUID, _ Priority) (Issue, error) {
	return Issue{}, errors.New("UpdateIssuePriority: not implemented")
}

// GetIssue retrieves a single issue by its ID.
func (s *IssueService) GetIssue(_ context.Context, _ uuid.UUID) (Issue, error) {
	return Issue{}, fmt.Errorf("get issue: not implemented")
}

// UpdateIssueStatus updates the status of the issue with the given ID.
func (s *IssueService) UpdateIssueStatus(_ context.Context, _ uuid.UUID, _ Status) (Issue, error) {
	return Issue{}, errors.New("not implemented")
}
