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

	expected := issuedocument.IssueDocument{
		ID:          id,
		Title:       title,
		Description: description,
		UpdatedAt:   updatedAt,
	}

	actual := issuedocument.NewIssueDocument(id, title, description, updatedAt)

	require.Equal(t, expected, actual)
}
