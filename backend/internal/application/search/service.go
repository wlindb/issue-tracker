package search

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/wlindb/issue-tracker/internal/application/api/model"
	"github.com/wlindb/issue-tracker/internal/domain/search/issuedocument"
)

type Searcher struct {
	issueDocumentRepository issuedocument.IssueDocumentRepository
	issueRepository         IssueRepository
}

var _ SearchService = (*Searcher)(nil)

type IssueRepository interface {
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Issue, error)
}

func NewSearcher(
	issueDocumentRepository issuedocument.IssueDocumentRepository,
	issueRepository IssueRepository,
) *Searcher {
	return &Searcher{
		issueDocumentRepository: issueDocumentRepository,
		issueRepository:         issueRepository,
	}
}

func (service Searcher) SearchIssues(ctx context.Context, description string) ([]model.Issue, error) {
	var zero []model.Issue
	documents, err := service.issueDocumentRepository.Find(ctx, description)
	if err != nil {
		return zero, fmt.Errorf("find: %w", err)
	}

	issues, err := service.issueRepository.GetByIDs(ctx, documents.IDs())
	if err != nil {
		return zero, fmt.Errorf("GetByIDs: %w", err)
	}

	return issues, nil
}
