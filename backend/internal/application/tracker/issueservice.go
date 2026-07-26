package tracker

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/wlindb/issue-tracker/internal/application/api/model"
	"github.com/wlindb/issue-tracker/internal/application/tracker/api"
	"github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
)

type IssueRepository interface {
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]issue.Issue, error)
}

type IssueService struct {
	repository IssueRepository
}

func NewIssueService(repository IssueRepository) IssueService {
	return IssueService{repository: repository}
}

func (s *IssueService) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Issue, error) {
	issues, err := s.repository.GetByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("GetByIDs: %w", err)
	}

	return api.IssuesFromDomain(issues), nil
}
