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

// UpdateTitle updates the title of the issue document identified by id.
// If the document's title already equals the given title the call is a no-op,
// which makes it safe to retry on at-least-once event redelivery.
func (s *IssueDocumentService) UpdateTitle(ctx context.Context, id uuid.UUID, title string) error {
	current, err := s.repository.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("update issue document title: %w", err)
	}
	if current.Title == title {
		return nil
	}

	updated := current
	updated.Title = title
	if _, err := s.repository.Update(ctx, updated); err != nil {
		return fmt.Errorf("update issue document title: %w", err)
	}
	return nil
}

// UpdateDescription updates the description of the issue document identified by id.
// If the document's description already equals the given description the call is a no-op,
// which makes it safe to retry on at-least-once event redelivery.
func (s *IssueDocumentService) UpdateDescription(ctx context.Context, id uuid.UUID, description string) error {
	current, err := s.repository.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("update issue document description: %w", err)
	}
	if current.Description == description {
		return nil
	}

	updated := current
	updated.Description = description
	if _, err := s.repository.Update(ctx, updated); err != nil {
		return fmt.Errorf("update issue document description: %w", err)
	}
	return nil
}
