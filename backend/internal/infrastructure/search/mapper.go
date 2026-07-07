package search

import (
	issuedocumentdomain "github.com/wlindb/issue-tracker/internal/domain/search/issuedocument"
	searchdb "github.com/wlindb/issue-tracker/internal/infrastructure/search/generated"
)

func issueDocumentToDomain(row searchdb.IssueDocument) issuedocumentdomain.IssueDocument {
	return issuedocumentdomain.IssueDocument{
		ID:          row.ID,
		WorkspaceID: row.WorkspaceID,
		Title:       row.Title,
		Description: row.Description,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}
}

func issueDocumentsToDomain(rows []searchdb.IssueDocument) []issuedocumentdomain.IssueDocument {
	documents := make([]issuedocumentdomain.IssueDocument, len(rows))
	for i, row := range rows {
		documents[i] = issueDocumentToDomain(row)
	}
	return documents
}

func createIssueDocumentParamsFromDomain(document issuedocumentdomain.IssueDocument) searchdb.CreateIssueDocumentParams {
	return searchdb.CreateIssueDocumentParams{
		ID:          document.ID,
		WorkspaceID: document.WorkspaceID,
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
