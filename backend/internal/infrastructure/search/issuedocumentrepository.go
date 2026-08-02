package search

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pgvector/pgvector-go"

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

// Update persists changes to an issue document using optimistic locking on updated_at.
// Returns ErrUpdateConflict if the document was modified, or no longer exists, since it was read.
func (r *IssueDocumentRepository) Update(ctx context.Context, document issuedocumentdomain.IssueDocument) (issuedocumentdomain.IssueDocument, error) {
	queries := searchdb.New(r.db)
	row, err := queries.UpdateIssueDocument(ctx, updateIssueDocumentParamsFromDomain(document))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return issuedocumentdomain.IssueDocument{}, fmt.Errorf("update issue document: %w", issuedocumentdomain.ErrUpdateConflict)
		}
		return issuedocumentdomain.IssueDocument{}, fmt.Errorf("update issue document: %w", err)
	}
	return issueDocumentToDomain(row), nil
}

// Get retrieves a single issue document by its ID.
// Returns ErrIssueDocumentNotFound if no document with the given ID exists.
func (r *IssueDocumentRepository) Get(ctx context.Context, id uuid.UUID) (issuedocumentdomain.IssueDocument, error) {
	queries := searchdb.New(r.db)
	row, err := queries.GetIssueDocument(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return issuedocumentdomain.IssueDocument{}, fmt.Errorf("get issue document: %w", issuedocumentdomain.ErrIssueDocumentNotFound)
		}
		return issuedocumentdomain.IssueDocument{}, fmt.Errorf("get issue document: %w", err)
	}
	return issueDocumentToDomain(row), nil
}

func (r *IssueDocumentRepository) Find(ctx context.Context, description string, embedding []float32) (issuedocumentdomain.IssueDocuments, error) {
	queries := searchdb.New(r.db)
	rows, err := queries.FindIssueDocumentsByDescription(ctx, searchdb.FindIssueDocumentsByDescriptionParams{
		Description: description,
		Embedding:   pgvector.NewVector(embedding),
	})
	if err != nil {
		return issuedocumentdomain.IssueDocuments{}, fmt.Errorf("find issue documents by description: %w", err)
	}
	if len(rows) == 0 {
		return issuedocumentdomain.IssueDocuments{}, nil
	}
	return issueDocumentSearchRowsToDomain(rows), nil
}
