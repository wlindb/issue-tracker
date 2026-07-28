// Package issuedocument
package issuedocument

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type IssueDocuments struct {
	Documents []IssueDocument
}

func (documents IssueDocuments) IDs() []uuid.UUID {
	ids := make([]uuid.UUID, len(documents.Documents))
	for i, id := range documents.Documents {
		ids[i] = id.ID
	}
	return ids
}

type IssueDocument struct {
	ID          uuid.UUID
	Title       string
	Description string
	UpdatedAt   time.Time
}

// IssueDocumentRepository defines the persistence interface for issue documents.
type IssueDocumentRepository interface {
	Create(ctx context.Context, document IssueDocument) (IssueDocument, error)
	Update(ctx context.Context, document IssueDocument) (IssueDocument, error)
	Find(ctx context.Context, description string) (IssueDocuments, error)
}
