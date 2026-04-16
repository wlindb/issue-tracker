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

// Create persists a comment and returns the stored result.
func (s *CommentService) Create(ctx context.Context, comment Comment) (*Comment, error) {
	result, err := s.repository.Create(ctx, comment)
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
