package issuedocument

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type IssueDocumentService struct {
	repository IssueDocumentRepository
}

func NewIssueDocumentService(repository IssueDocumentRepository) *IssueDocumentService {
	return &IssueDocumentService{repository: repository}
}

func (s *IssueDocumentService) Create(ctx context.Context, id uuid.UUID, title string, description string, updatedAt time.Time) error {
	document := NewIssueDocument(id, title, description, updatedAt)

	if _, err := s.repository.Create(ctx, document); err != nil {
		return fmt.Errorf("create issue document: %w", err)
	}
	return nil
}
