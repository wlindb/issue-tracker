package search

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	issuedocumentdomain "github.com/wlindb/issue-tracker/internal/domain/search/issuedocument"
	searchdb "github.com/wlindb/issue-tracker/internal/infrastructure/search/generated"
)

var _ issuedocumentdomain.IssueDocumentRepository = (*IssueDocumentRepository)(nil)

type IssueDocumentRepository struct {
	db searchdb.DBTX
}

func NewIssueDocumentRepository(db searchdb.DBTX) *IssueDocumentRepository {
	return &IssueDocumentRepository{db: db}
}

func (r *IssueDocumentRepository) Create(ctx context.Context, document issuedocumentdomain.IssueDocument) (issuedocumentdomain.IssueDocument, error) {
	queries := searchdb.New(r.db)
	row, err := queries.CreateIssueDocument(ctx, createIssueDocumentParamsFromDomain(document))
	if err == nil {
		return issueDocumentToDomain(row), nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return issuedocumentdomain.IssueDocument{}, fmt.Errorf("create issue document: %w", err)
	}

	existing, err := queries.GetIssueDocument(ctx, document.ID)
	if err != nil {
		return issuedocumentdomain.IssueDocument{}, fmt.Errorf("create issue document: get existing after conflict: %w", err)
	}
	return issueDocumentToDomain(existing), nil
}

func (r *IssueDocumentRepository) Update(ctx context.Context, document issuedocumentdomain.IssueDocument) (issuedocumentdomain.IssueDocument, error) {
	queries := searchdb.New(r.db)
	row, err := queries.UpdateIssueDocument(ctx, updateIssueDocumentParamsFromDomain(document))
	if err != nil {
		return issuedocumentdomain.IssueDocument{}, fmt.Errorf("update issue document: %w", err)
	}
	return issueDocumentToDomain(row), nil
}

func (r *IssueDocumentRepository) Find(ctx context.Context, description string) (issuedocumentdomain.IssueDocuments, error) {
	queries := searchdb.New(r.db)
	rows, err := queries.FindIssueDocumentsByDescription(ctx, description)
	if err != nil {
		return issuedocumentdomain.IssueDocuments{}, fmt.Errorf("find issue documents by description: %w", err)
	}
	if len(rows) == 0 {
		return issuedocumentdomain.IssueDocuments{}, nil
	}
	return issueDocumentsToDomain(rows), nil
}
