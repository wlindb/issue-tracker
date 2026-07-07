//go:build !integration

package search

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	issuedocumentdomain "github.com/wlindb/issue-tracker/internal/domain/search/issuedocument"
	searchdb "github.com/wlindb/issue-tracker/internal/infrastructure/search/generated"
)

func Test_issueDocumentToDomain_Row_ReturnsDomainIssueDocument(t *testing.T) {
	id := uuid.New()
	workspaceID := uuid.New()
	now := time.Now().UTC()
	row := searchdb.IssueDocument{
		ID:          id,
		WorkspaceID: workspaceID,
		Title:       "Fix bug",
		Description: "detailed description",
		CreatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
	}

	actual := issueDocumentToDomain(row)

	assert.Equal(t, id, actual.ID)
	assert.Equal(t, workspaceID, actual.WorkspaceID)
	assert.Equal(t, "Fix bug", actual.Title)
	assert.Equal(t, "detailed description", actual.Description)
	assert.Equal(t, now, actual.CreatedAt)
	assert.Equal(t, now, actual.UpdatedAt)
}

func Test_issueDocumentsToDomain_Empty_ReturnsEmptySlice(t *testing.T) {
	actual := issueDocumentsToDomain([]searchdb.IssueDocument{})

	assert.NotNil(t, actual)
	assert.Empty(t, actual)
}

func Test_issueDocumentsToDomain_MultipleRows_ReturnsMappedIssueDocuments(t *testing.T) {
	now := time.Now().UTC()
	firstID, secondID := uuid.New(), uuid.New()
	rows := []searchdb.IssueDocument{
		{
			ID:          firstID,
			WorkspaceID: uuid.New(),
			Title:       "First",
			Description: "First description",
			CreatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
			UpdatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
		},
		{
			ID:          secondID,
			WorkspaceID: uuid.New(),
			Title:       "Second",
			Description: "Second description",
			CreatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
			UpdatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
		},
	}

	actual := issueDocumentsToDomain(rows)

	require.Len(t, actual, 2)
	assert.Equal(t, firstID, actual[0].ID)
	assert.Equal(t, "First", actual[0].Title)
	assert.Equal(t, secondID, actual[1].ID)
	assert.Equal(t, "Second", actual[1].Title)
}

func Test_createIssueDocumentParamsFromDomain_IssueDocument_ReturnsParams(t *testing.T) {
	id := uuid.New()
	workspaceID := uuid.New()
	domainDocument := issuedocumentdomain.IssueDocument{
		ID:          id,
		WorkspaceID: workspaceID,
		Title:       "Title",
		Description: "Description",
	}

	actual := createIssueDocumentParamsFromDomain(domainDocument)

	assert.Equal(t, id, actual.ID)
	assert.Equal(t, workspaceID, actual.WorkspaceID)
	assert.Equal(t, "Title", actual.Title)
	assert.Equal(t, "Description", actual.Description)
}

func Test_updateIssueDocumentParamsFromDomain_IssueDocument_ReturnsParams(t *testing.T) {
	id := uuid.New()
	domainDocument := issuedocumentdomain.IssueDocument{
		ID:          id,
		WorkspaceID: uuid.New(),
		Title:       "Updated title",
		Description: "Updated description",
	}

	actual := updateIssueDocumentParamsFromDomain(domainDocument)

	assert.Equal(t, id, actual.ID)
	assert.Equal(t, "Updated title", actual.Title)
	assert.Equal(t, "Updated description", actual.Description)
}
