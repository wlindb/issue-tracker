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

func Test_Create_ValidDocument_SuccessfulIssueDocumentCreation(t *testing.T) {
	repository := search.NewIssueDocumentRepository(testPool)
	id := uuid.New()
	workspaceID := uuid.New()

	actual, err := repository.Create(context.Background(), issuedocumentdomain.IssueDocument{
		ID:          id,
		WorkspaceID: workspaceID,
		Title:       "Fix bug",
		Description: "A detailed bug description",
	})

	require.NoError(t, err)
	assert.Equal(t, id, actual.ID)
	assert.Equal(t, workspaceID, actual.WorkspaceID)
	assert.Equal(t, "Fix bug", actual.Title)
	assert.Equal(t, "A detailed bug description", actual.Description)
	assert.False(t, actual.CreatedAt.IsZero())
	assert.False(t, actual.UpdatedAt.IsZero())
}

func Test_Create_DuplicateID_ReturnsError(t *testing.T) {
	repository := search.NewIssueDocumentRepository(testPool)
	id := uuid.New()
	document := issuedocumentdomain.IssueDocument{
		ID:          id,
		WorkspaceID: uuid.New(),
		Title:       "First",
		Description: "First description",
	}

	_, err := repository.Create(context.Background(), document)
	require.NoError(t, err)

	_, err = repository.Create(context.Background(), document)

	require.Error(t, err)
}

func Test_Update_ExistingDocument_SuccessfulUpdate(t *testing.T) {
	repository := search.NewIssueDocumentRepository(testPool)
	created, err := repository.Create(context.Background(), issuedocumentdomain.IssueDocument{
		ID:          uuid.New(),
		WorkspaceID: uuid.New(),
		Title:       "Original title",
		Description: "Original description",
	})
	require.NoError(t, err)

	created.Title = "Updated title"
	created.Description = "Updated description"

	actual, err := repository.Update(context.Background(), created)

	require.NoError(t, err)
	assert.Equal(t, "Updated title", actual.Title)
	assert.Equal(t, "Updated description", actual.Description)
	assert.True(t, actual.UpdatedAt.After(created.CreatedAt) || actual.UpdatedAt.Equal(created.CreatedAt))
}

func Test_Update_NonExistentDocument_ReturnsError(t *testing.T) {
	repository := search.NewIssueDocumentRepository(testPool)

	_, err := repository.Update(context.Background(), issuedocumentdomain.IssueDocument{
		ID:          uuid.New(),
		Title:       "Missing",
		Description: "Missing description",
	})

	require.Error(t, err)
}

func Test_Find_MatchingDescription_ReturnsIssueDocuments(t *testing.T) {
	repository := search.NewIssueDocumentRepository(testPool)
	workspaceID := uuid.New()

	_, err := repository.Create(context.Background(), issuedocumentdomain.IssueDocument{
		ID:          uuid.New(),
		WorkspaceID: workspaceID,
		Title:       "Backend issue",
		Description: "database connection pooling bug",
	})
	require.NoError(t, err)

	actual, err := repository.Find(context.Background(), "database")

	require.NoError(t, err)
	assert.NotEmpty(t, actual)
}

func Test_Find_NoMatch_ReturnsEmptySlice(t *testing.T) {
	repository := search.NewIssueDocumentRepository(testPool)

	actual, err := repository.Find(context.Background(), "zzz-no-match-"+uuid.New().String())

	require.NoError(t, err)
	assert.Empty(t, actual)
}
