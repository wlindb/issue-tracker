package issue

import (
	"context"
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

// CreateIssue creates a new issue from the given command.
func (s *IssueService) CreateIssue(ctx context.Context, command CreateIssueCommand) (*Issue, error) {
	issue := command.ToIssue(uuid.New(), command.Slugify)
	result, err := s.repository.CreateIssue(ctx, issue)
	if err != nil {
		return nil, fmt.Errorf("create issue: %w", err)
	}
	return result, nil
}
