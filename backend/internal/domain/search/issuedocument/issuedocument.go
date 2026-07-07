package issuedocument

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// IssueDocument is a search-indexed representation of an issue.
type IssueDocument struct {
	ID          uuid.UUID
	WorkspaceID uuid.UUID
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// IssueDocumentRepository defines the persistence interface for issue documents.
type IssueDocumentRepository interface {
	Create(ctx context.Context, document IssueDocument) (IssueDocument, error)
	Update(ctx context.Context, document IssueDocument) (IssueDocument, error)
	Find(ctx context.Context, description string) ([]IssueDocument, error)
}
