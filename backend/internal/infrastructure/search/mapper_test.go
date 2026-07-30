//go:build !integration

package search

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pgvector/pgvector-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	issuedocumentdomain "github.com/wlindb/issue-tracker/internal/domain/search/issuedocument"
	searchdb "github.com/wlindb/issue-tracker/internal/infrastructure/search/generated"
)

func Test_issueDocumentToDomain_Row_ReturnsDomainIssueDocument(t *testing.T) {
	id := uuid.New()
	workspaceID := uuid.New()
	now := time.Now().UTC()
	embedding := []float32{0.1, 0.2, 0.3}
	vector := pgvector.NewVector(embedding)
	row := searchdb.IssueDocument{
		ID:          id,
		WorkspaceID: workspaceID,
		Title:       "Fix bug",
		Description: "detailed description",
		CreatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
		Embedding:   &vector,
	}

	actual := issueDocumentToDomain(row)

	assert.Equal(t, id, actual.ID)
	assert.Equal(t, "Fix bug", actual.Title)
	assert.Equal(t, "detailed description", actual.Description)
	assert.Equal(t, now, actual.UpdatedAt)
	assert.Equal(t, embedding, actual.Embedding)
}

func Test_issueDocumentToDomain_NilEmbedding_ReturnsNilEmbedding(t *testing.T) {
	row := searchdb.IssueDocument{
		ID:          uuid.New(),
		WorkspaceID: uuid.New(),
		Title:       "Legacy issue",
		Description: "created before embeddings existed",
		CreatedAt:   pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true},
		UpdatedAt:   pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true},
		Embedding:   nil,
	}

	actual := issueDocumentToDomain(row)

	assert.Nil(t, actual.Embedding)
}

func Test_issueDocumentsToDomain_Empty_ReturnsEmptySlice(t *testing.T) {
	actual := issueDocumentsToDomain([]searchdb.IssueDocument{})

	assert.NotNil(t, actual)
	assert.Empty(t, actual.Documents)
}

func Test_issueDocumentsToDomain_MultipleRows_ReturnsMappedIssueDocuments(t *testing.T) {
	now := time.Now().UTC()
	firstID, secondID := uuid.New(), uuid.New()
	firstEmbedding := []float32{0.1, 0.2, 0.3}
	secondEmbedding := []float32{0.4, 0.5, 0.6}
	firstVector := pgvector.NewVector(firstEmbedding)
	secondVector := pgvector.NewVector(secondEmbedding)
	rows := []searchdb.IssueDocument{
		{
			ID:          firstID,
			WorkspaceID: uuid.New(),
			Title:       "First",
			Description: "First description",
			CreatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
			UpdatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
			Embedding:   &firstVector,
		},
		{
			ID:          secondID,
			WorkspaceID: uuid.New(),
			Title:       "Second",
			Description: "Second description",
			CreatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
			UpdatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
			Embedding:   &secondVector,
		},
	}

	actual := issueDocumentsToDomain(rows)

	require.Len(t, actual.Documents, 2)
	assert.Equal(t, firstID, actual.Documents[0].ID)
	assert.Equal(t, "First", actual.Documents[0].Title)
	assert.Equal(t, firstEmbedding, actual.Documents[0].Embedding)
	assert.Equal(t, secondID, actual.Documents[1].ID)
	assert.Equal(t, "Second", actual.Documents[1].Title)
	assert.Equal(t, secondEmbedding, actual.Documents[1].Embedding)
}

func Test_createIssueDocumentParamsFromDomain_IssueDocument_ReturnsParams(t *testing.T) {
	id := uuid.New()
	embedding := []float32{0.1, 0.2, 0.3}
	domainDocument := issuedocumentdomain.IssueDocument{
		ID:          id,
		Title:       "Title",
		Description: "Description",
		Embedding:   embedding,
	}

	actual := createIssueDocumentParamsFromDomain(domainDocument)

	assert.Equal(t, id, actual.ID)
	assert.Equal(t, "Title", actual.Title)
	assert.Equal(t, "Description", actual.Description)
	assert.Equal(t, embedding, actual.Embedding.Slice())
}

func Test_updateIssueDocumentParamsFromDomain_IssueDocument_ReturnsParams(t *testing.T) {
	id := uuid.New()
	now := time.Now().UTC()
	domainDocument := issuedocumentdomain.IssueDocument{
		ID:          id,
		Title:       "Updated title",
		Description: "Updated description",
		UpdatedAt:   now,
	}

	actual := updateIssueDocumentParamsFromDomain(domainDocument)

	assert.Equal(t, id, actual.ID)
	assert.Equal(t, "Updated title", actual.Title)
	assert.Equal(t, "Updated description", actual.Description)
	assert.Equal(t, now, actual.UpdatedAt.Time)
	assert.True(t, actual.UpdatedAt.Valid)
}
