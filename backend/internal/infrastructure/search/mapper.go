package search

import (
	issuedocumentdomain "github.com/wlindb/issue-tracker/internal/domain/search/issuedocument"
	searchdb "github.com/wlindb/issue-tracker/internal/infrastructure/search/generated"
)

func issueDocumentToDomain(row searchdb.IssueDocument) issuedocumentdomain.IssueDocument {
	return issuedocumentdomain.IssueDocument{
		ID:          row.ID,
		Title:       row.Title,
		Description: row.Description,
		UpdatedAt:   row.UpdatedAt.Time,
	}
}

func issueDocumentsToDomain(rows []searchdb.IssueDocument) issuedocumentdomain.IssueDocuments {
	documents := make([]issuedocumentdomain.IssueDocument, len(rows))
	for i, row := range rows {
		documents[i] = issueDocumentToDomain(row)
	}
	return issuedocumentdomain.IssueDocuments{Documents: documents}
}

func createIssueDocumentParamsFromDomain(document issuedocumentdomain.IssueDocument) searchdb.CreateIssueDocumentParams {
	return searchdb.CreateIssueDocumentParams{
		ID:          document.ID,
		Title:       document.Title,
		Description: document.Description,
	}
}

func updateIssueDocumentParamsFromDomain(document issuedocumentdomain.IssueDocument) searchdb.UpdateIssueDocumentParams {
	return searchdb.UpdateIssueDocumentParams{
		ID:          document.ID,
		Title:       document.Title,
		Description: document.Description,
	}
}
