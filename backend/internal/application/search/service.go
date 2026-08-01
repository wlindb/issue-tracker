package search

import (
	"context"
	"fmt"
	"slices"

	"github.com/google/uuid"

	"github.com/wlindb/issue-tracker/internal/application/api/model"
	"github.com/wlindb/issue-tracker/internal/domain/search/issuedocument"
)

type Searcher struct {
	issueDocumentFinder IssueDocumentFinder
	issueRepository     IssueRepository
}

var _ SearchService = (*Searcher)(nil)

type IssueRepository interface {
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Issue, error)
}

type IssueDocumentFinder interface {
	Find(ctx context.Context, description string) (issuedocument.IssueDocuments, error)
}

func NewSearcher(
	issueDocumentFinder IssueDocumentFinder,
	issueRepository IssueRepository,
) *Searcher {
	return &Searcher{
		issueDocumentFinder: issueDocumentFinder,
		issueRepository:     issueRepository,
	}
}

func (service Searcher) SearchIssues(ctx context.Context, description string) ([]model.Issue, error) {
	var zero []model.Issue
	documents, err := service.issueDocumentFinder.Find(ctx, description)
	if err != nil {
		return zero, fmt.Errorf("find: %w", err)
	}

	issues, err := service.issueRepository.GetByIDs(ctx, documents.IDs())
	if err != nil {
		return zero, fmt.Errorf("GetByIDs: %w", err)
	}
	slices.SortFunc(issues, func(a, b model.Issue) int {
		return documents.CompareRank(a.Id, b.Id)
	})

	return issues, nil
}
