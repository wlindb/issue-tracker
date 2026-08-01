package issuedocument_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/wlindb/issue-tracker/internal/domain/search/issuedocument"
)

func Test_NewIssueDocument_ValidInput_SuccessfulIssueDocumentCreation(t *testing.T) {
	id := uuid.New()
	title := "test issue"
	description := "test description"
	updatedAt := time.Now()
	embedding := []float32{0.1, 0.2, 0.3}

	expected := issuedocument.IssueDocument{
		ID:          id,
		Title:       title,
		Description: description,
		UpdatedAt:   updatedAt,
		Embedding:   embedding,
	}

	actual := issuedocument.NewIssueDocument(id, title, description, updatedAt, embedding)

	require.Equal(t, expected, actual)
}

func Test_CompareRank_KnownOrder_ReturnsExpectedSign(t *testing.T) {
	firstID, secondID, thirdID := uuid.New(), uuid.New(), uuid.New()
	documents := issuedocument.NewIssueDocuments([]issuedocument.IssueDocument{
		{ID: firstID},
		{ID: secondID},
		{ID: thirdID},
	})

	require.Negative(t, documents.CompareRank(firstID, secondID))
	require.Positive(t, documents.CompareRank(thirdID, secondID))
	require.Zero(t, documents.CompareRank(secondID, secondID))
}

func Test_CompareRank_UnknownID_TreatsAsRankZero(t *testing.T) {
	firstID, secondID := uuid.New(), uuid.New()
	unknownID := uuid.New()
	documents := issuedocument.NewIssueDocuments([]issuedocument.IssueDocument{
		{ID: firstID},
		{ID: secondID},
	})

	// unknownID defaults to rank 0, same as firstID, so it only differs from secondID (rank 1).
	require.Zero(t, documents.CompareRank(firstID, unknownID))
	require.Positive(t, documents.CompareRank(secondID, unknownID))
}
