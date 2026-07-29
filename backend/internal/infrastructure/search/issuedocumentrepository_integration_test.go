//go:build integration

package search_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	issuedocumentdomain "github.com/wlindb/issue-tracker/internal/domain/search/issuedocument"
	"github.com/wlindb/issue-tracker/internal/infrastructure/search"
)

func Test_Create_ValidDocument_SuccessfulIssueDocumentCreation(t *testing.T) {
	repository := search.NewIssueDocumentRepository(testPool)
	id := uuid.New()
	workspaceID := uuid.New()
	ctx := withWorkspaceContext(workspaceID, uuid.New())

	actual, err := repository.Create(ctx, issuedocumentdomain.IssueDocument{
		ID:          id,
		Title:       "Fix bug",
		Description: "A detailed bug description",
	})

	require.NoError(t, err)
	assert.Equal(t, id, actual.ID)
	assert.Equal(t, "Fix bug", actual.Title)
	assert.Equal(t, "A detailed bug description", actual.Description)
	assert.False(t, actual.UpdatedAt.IsZero())
}

func Test_Create_DuplicateID_ReturnsExistingDocumentNoError(t *testing.T) {
	repository := search.NewIssueDocumentRepository(testPool)
	id := uuid.New()
	ctx := withWorkspaceContext(uuid.New(), uuid.New())

	created, err := repository.Create(ctx, issuedocumentdomain.IssueDocument{
		ID:          id,
		Title:       "First",
		Description: "First description",
	})
	require.NoError(t, err)

	actual, err := repository.Create(ctx, issuedocumentdomain.IssueDocument{
		ID:          id,
		Title:       "Second",
		Description: "Second description",
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

func Test_Update_NonExistentDocument_ReturnsError(t *testing.T) {
	repository := search.NewIssueDocumentRepository(testPool)

	_, err := repository.Update(withWorkspaceContext(uuid.New(), uuid.New()), issuedocumentdomain.IssueDocument{
		ID:          uuid.New(),
		Title:       "Missing",
		Description: "Missing description",
	})

	require.Error(t, err)
}

func Test_Find_MatchingDescription_ReturnsIssueDocuments(t *testing.T) {
	repository := search.NewIssueDocumentRepository(testPool)
	ctx := withWorkspaceContext(uuid.New(), uuid.New())

	_, err := repository.Create(ctx, issuedocumentdomain.IssueDocument{
		ID:          uuid.New(),
		Title:       "Backend issue",
		Description: "database connection pooling bug",
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
