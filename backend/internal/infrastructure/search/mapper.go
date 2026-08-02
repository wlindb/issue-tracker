package search

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pgvector/pgvector-go"

	issuedocumentdomain "github.com/wlindb/issue-tracker/internal/domain/search/issuedocument"
	searchdb "github.com/wlindb/issue-tracker/internal/infrastructure/search/generated"
)

func issueDocumentToDomain(row searchdb.IssueDocument) issuedocumentdomain.IssueDocument {
	var embedding []float32
	if row.Embedding != nil {
		embedding = row.Embedding.Slice()
	}
	return issuedocumentdomain.IssueDocument{
		ID:          row.ID,
		Title:       row.Title,
		Description: row.Description,
		UpdatedAt:   row.UpdatedAt.Time,
		Embedding:   embedding,
	}
}

func issueDocumentsToDomain(rows []searchdb.IssueDocument) issuedocumentdomain.IssueDocuments {
	documents := make([]issuedocumentdomain.IssueDocument, len(rows))
	for i, row := range rows {
		documents[i] = issueDocumentToDomain(row)
	}
	return issuedocumentdomain.NewIssueDocuments(documents)
}

func issueDocumentSearchRowToDomain(row searchdb.FindIssueDocumentsByDescriptionRow) issuedocumentdomain.IssueDocument {
	var embedding []float32
	if row.Embedding != nil {
		embedding = row.Embedding.Slice()
	}
	return issuedocumentdomain.IssueDocument{
		ID:          row.ID,
		Title:       row.Title,
		Description: row.Description,
		UpdatedAt:   row.UpdatedAt.Time,
		Embedding:   embedding,
	}
}

func issueDocumentSearchRowsToDomain(rows []searchdb.FindIssueDocumentsByDescriptionRow) issuedocumentdomain.IssueDocuments {
	documents := make([]issuedocumentdomain.IssueDocument, len(rows))
	for i, row := range rows {
		documents[i] = issueDocumentSearchRowToDomain(row)
	}
	return issuedocumentdomain.NewIssueDocuments(documents)
}

func createIssueDocumentParamsFromDomain(document issuedocumentdomain.IssueDocument) searchdb.CreateIssueDocumentParams {
	embedding := pgvector.NewVector(document.Embedding)
	return searchdb.CreateIssueDocumentParams{
		ID:          document.ID,
		Title:       document.Title,
		Description: document.Description,
		Embedding:   &embedding,
	}
}

func updateIssueDocumentParamsFromDomain(document issuedocumentdomain.IssueDocument) searchdb.UpdateIssueDocumentParams {
	return searchdb.UpdateIssueDocumentParams{
		ID:          document.ID,
		Title:       document.Title,
		Description: document.Description,
		UpdatedAt:   pgtype.Timestamptz{Time: document.UpdatedAt, Valid: true},
	}
}
