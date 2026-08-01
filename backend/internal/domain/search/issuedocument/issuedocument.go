// Package issuedocument
package issuedocument

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUpdateConflict        = errors.New("update conflict")
	ErrIssueDocumentNotFound = errors.New("issue document not found")
)

type IssueDocuments struct {
	Documents []IssueDocument
	rank      map[uuid.UUID]int
}

func NewIssueDocuments(documents []IssueDocument) IssueDocuments {
	rank := make(map[uuid.UUID]int, len(documents))
	for i, document := range documents {
		rank[document.ID] = i
	}
	return IssueDocuments{Documents: documents, rank: rank}
}

func (documents IssueDocuments) IDs() []uuid.UUID {
	ids := make([]uuid.UUID, len(documents.Documents))
	for i, id := range documents.Documents {
		ids[i] = id.ID
	}
	return ids
}

// CompareRank reports the relative rank order of two issue document IDs within this result
// set, suitable as a slices.SortFunc comparator: negative if firstID is a stronger match
// (ranks before) secondID, positive if weaker, zero if equal or either ID is not present.
func (documents IssueDocuments) CompareRank(firstID, secondID uuid.UUID) int {
	return documents.rank[firstID] - documents.rank[secondID]
}

type IssueDocument struct {
	ID          uuid.UUID
	Title       string
	Description string
	UpdatedAt   time.Time
	Embedding   []float32
}

func NewIssueDocument(id uuid.UUID, title string, description string, updatedAt time.Time, embedding []float32) IssueDocument {
	return IssueDocument{
		ID:          id,
		Title:       title,
		Description: description,
		UpdatedAt:   updatedAt,
		Embedding:   embedding,
	}
}

// IssueDocumentRepository defines the persistence interface for issue documents.
type IssueDocumentRepository interface {
	Create(ctx context.Context, document IssueDocument) (IssueDocument, error)
	Update(ctx context.Context, document IssueDocument) (IssueDocument, error)
	Get(ctx context.Context, id uuid.UUID) (IssueDocument, error)
	Find(ctx context.Context, description string, embedding []float32) (IssueDocuments, error)
}
