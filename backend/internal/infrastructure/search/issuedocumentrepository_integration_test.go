//go:build integration

package search_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	issuedocumentdomain "github.com/wlindb/issue-tracker/internal/domain/search/issuedocument"
	"github.com/wlindb/issue-tracker/internal/infrastructure/search"
)

// testEmbedding returns a 1536-dimensional vector matching the issue_documents.embedding column,
// so integration tests can round-trip a valid embedding without calling the real Gemini API.
func testEmbedding() []float32 {
	embedding := make([]float32, 1536)
	for i := range embedding {
		embedding[i] = 0.1
	}
	return embedding
}

func Test_Create_ValidDocument_SuccessfulIssueDocumentCreation(t *testing.T) {
	repository := search.NewIssueDocumentRepository(testPool)
	id := uuid.New()
	workspaceID := uuid.New()
	ctx := withWorkspaceContext(workspaceID, uuid.New())

	actual, err := repository.Create(ctx, issuedocumentdomain.IssueDocument{
		ID:          id,
		Title:       "Fix bug",
		Description: "A detailed bug description",
		Embedding:   testEmbedding(),
	})

	require.NoError(t, err)
	assert.Equal(t, id, actual.ID)
	assert.Equal(t, "Fix bug", actual.Title)
	assert.Equal(t, "A detailed bug description", actual.Description)
	assert.False(t, actual.UpdatedAt.IsZero())
	assert.Equal(t, testEmbedding(), actual.Embedding)
}

func Test_Create_DuplicateID_ReturnsExistingDocumentNoError(t *testing.T) {
	repository := search.NewIssueDocumentRepository(testPool)
	id := uuid.New()
	ctx := withWorkspaceContext(uuid.New(), uuid.New())

	created, err := repository.Create(ctx, issuedocumentdomain.IssueDocument{
		ID:          id,
		Title:       "First",
		Description: "First description",
		Embedding:   testEmbedding(),
	})
	require.NoError(t, err)

	actual, err := repository.Create(ctx, issuedocumentdomain.IssueDocument{
		ID:          id,
		Title:       "Second",
		Description: "Second description",
		Embedding:   testEmbedding(),
	})

	require.NoError(t, err)
	assert.Equal(t, id, actual.ID)
	assert.Equal(t, created.Title, actual.Title)
	assert.Equal(t, created.Description, actual.Description)
	assert.Equal(t, created.UpdatedAt, actual.UpdatedAt)
}

func Test_Update_ExistingDocument_SuccessfulUpdate(t *testing.T) {
	repository := search.NewIssueDocumentRepository(testPool)
	ctx := withWorkspaceContext(uuid.New(), uuid.New())
	created, err := repository.Create(ctx, issuedocumentdomain.IssueDocument{
		ID:          uuid.New(),
		Title:       "Original title",
		Description: "Original description",
		Embedding:   testEmbedding(),
	})
	require.NoError(t, err)

	created.Title = "Updated title"
	created.Description = "Updated description"

	actual, err := repository.Update(ctx, created)

	require.NoError(t, err)
	assert.Equal(t, "Updated title", actual.Title)
	assert.Equal(t, "Updated description", actual.Description)
	assert.True(t, actual.UpdatedAt.After(created.UpdatedAt) || actual.UpdatedAt.Equal(created.UpdatedAt))
}

func Test_Update_NonExistentDocument_ReturnsErrUpdateConflict(t *testing.T) {
	repository := search.NewIssueDocumentRepository(testPool)

	_, err := repository.Update(withWorkspaceContext(uuid.New(), uuid.New()), issuedocumentdomain.IssueDocument{
		ID:          uuid.New(),
		Title:       "Missing",
		Description: "Missing description",
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, issuedocumentdomain.ErrUpdateConflict)
}

func Test_Update_StaleEntity_ReturnsErrUpdateConflict(t *testing.T) {
	repository := search.NewIssueDocumentRepository(testPool)
	ctx := withWorkspaceContext(uuid.New(), uuid.New())
	created, err := repository.Create(ctx, issuedocumentdomain.IssueDocument{
		ID:          uuid.New(),
		Title:       "Original title",
		Description: "Original description",
		Embedding:   testEmbedding(),
	})
	require.NoError(t, err)

	firstUpdate := created
	firstUpdate.Title = "First update"
	_, err = repository.Update(ctx, firstUpdate)
	require.NoError(t, err)

	staleUpdate := created
	staleUpdate.Title = "Second update"
	_, err = repository.Update(ctx, staleUpdate)

	require.Error(t, err)
	assert.ErrorIs(t, err, issuedocumentdomain.ErrUpdateConflict)
}

func Test_Get_ExistingDocument_ReturnsDocument(t *testing.T) {
	repository := search.NewIssueDocumentRepository(testPool)
	ctx := withWorkspaceContext(uuid.New(), uuid.New())
	created, err := repository.Create(ctx, issuedocumentdomain.IssueDocument{
		ID:          uuid.New(),
		Title:       "Fix bug",
		Description: "A detailed bug description",
		Embedding:   testEmbedding(),
	})
	require.NoError(t, err)

	actual, err := repository.Get(ctx, created.ID)

	require.NoError(t, err)
	assert.Equal(t, created.ID, actual.ID)
	assert.Equal(t, "Fix bug", actual.Title)
}

func Test_Get_NonExistentDocument_ReturnsErrIssueDocumentNotFound(t *testing.T) {
	repository := search.NewIssueDocumentRepository(testPool)

	_, err := repository.Get(withWorkspaceContext(uuid.New(), uuid.New()), uuid.New())

	require.Error(t, err)
	assert.ErrorIs(t, err, issuedocumentdomain.ErrIssueDocumentNotFound)
}

// Test_Get_NullEmbedding_ReturnsDocumentWithNilEmbedding covers rows created before embedding-writing
// existed (or any row where embedding generation was skipped), where the embedding column is SQL NULL.
// Scanning NULL into a non-nullable pgvector.Vector used to panic; see mapper.go's issueDocumentToDomain.
func Test_Get_NullEmbedding_ReturnsDocumentWithNilEmbedding(t *testing.T) {
	ctx := withWorkspaceContext(uuid.New(), uuid.New())
	id := uuid.New()
	insertIssueDocumentWithoutEmbedding(t, ctx, id, "Legacy issue", "created before embeddings existed")

	repository := search.NewIssueDocumentRepository(testPool)
	actual, err := repository.Get(ctx, id)

	require.NoError(t, err)
	assert.Equal(t, id, actual.ID)
	assert.Nil(t, actual.Embedding)
}

func Test_Find_MatchingDescription_ReturnsIssueDocuments(t *testing.T) {
	repository := search.NewIssueDocumentRepository(testPool)
	ctx := withWorkspaceContext(uuid.New(), uuid.New())

	_, err := repository.Create(ctx, issuedocumentdomain.IssueDocument{
		ID:          uuid.New(),
		Title:       "Backend issue",
		Description: "database connection pooling bug",
		Embedding:   testEmbedding(),
	})
	require.NoError(t, err)

	actual, err := repository.Find(ctx, "database")

	require.NoError(t, err)
	assert.NotEmpty(t, actual.Documents)
}

func Test_Find_NoMatch_ReturnsEmptySlice(t *testing.T) {
	repository := search.NewIssueDocumentRepository(testPool)

	actual, err := repository.Find(withWorkspaceContext(uuid.New(), uuid.New()), "zzz-no-match-"+uuid.New().String())

	require.NoError(t, err)
	assert.Empty(t, actual)
}

// Test_Find_NullEmbedding_ReturnsIssueDocuments is a regression test for the panic reported when
// Find scanned a row whose embedding column is SQL NULL (see Test_Get_NullEmbedding_ReturnsDocumentWithNilEmbedding).
func Test_Find_NullEmbedding_ReturnsIssueDocuments(t *testing.T) {
	ctx := withWorkspaceContext(uuid.New(), uuid.New())
	id := uuid.New()
	insertIssueDocumentWithoutEmbedding(t, ctx, id, "Legacy issue", "no embedding was ever generated for this row")

	repository := search.NewIssueDocumentRepository(testPool)
	actual, err := repository.Find(ctx, "no embedding was ever generated")

	require.NoError(t, err)
	assert.Contains(t, actual.IDs(), id)
}

// insertIssueDocumentWithoutEmbedding inserts an issue_documents row leaving embedding as SQL NULL,
// bypassing IssueDocumentRepository.Create which always writes a non-null vector.
func insertIssueDocumentWithoutEmbedding(t *testing.T, ctx context.Context, id uuid.UUID, title, description string) {
	t.Helper()
	_, err := testPool.Exec(ctx, `
		INSERT INTO issue_documents (id, workspace_id, title, description, created_at, updated_at)
		VALUES ($1, current_setting('app.workspace_id')::uuid, $2, $3, NOW(), NOW())
	`, id, title, description)
	require.NoError(t, err)
}
