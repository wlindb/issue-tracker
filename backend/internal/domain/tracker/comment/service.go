package comment

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// CommentService implements the domain logic for managing comments.
type CommentService struct {
	repository Repository
}

// NewCommentService creates a CommentService wired to the given repository.
func NewCommentService(repository Repository) *CommentService {
	return &CommentService{repository: repository}
}

// Create creates a new comment and persists it.
func (s *CommentService) Create(ctx context.Context, id uuid.UUID, issueID uuid.UUID, authorID uuid.UUID, body string) (*Comment, error) {
	c, err := New(id, body, authorID, issueID)
	if err != nil {
		return nil, fmt.Errorf("create comment: %w", err)
	}
	result, err := s.repository.Create(ctx, c)
	if err != nil {
		return nil, fmt.Errorf("create comment: %w", err)
	}
	return &result, nil
}

// List returns all comments for the given issue.
func (s *CommentService) List(ctx context.Context, issueID uuid.UUID, _ ListCommentQuery) (Comments, error) {
	items, err := s.repository.Get(ctx, issueID)
	if err != nil {
		return Comments{}, fmt.Errorf("list comments: %w", err)
	}
	return Comments{Items: items}, nil
}
