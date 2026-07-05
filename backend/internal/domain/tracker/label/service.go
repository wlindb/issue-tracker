package label

import (
	"context"
	"fmt"
)

// LabelService implements the domain logic for managing labels.
type LabelService struct {
	repository LabelRepository
}

// NewLabelService creates a LabelService wired to the given repository.
func NewLabelService(repository LabelRepository) *LabelService {
	return &LabelService{repository: repository}
}

// Create gets or creates a label with the given name and returns it.
func (s *LabelService) Create(ctx context.Context, name string) (Label, error) {
	result, err := s.repository.GetOrCreate(ctx, name)
	if err != nil {
		return Label{}, fmt.Errorf("create label: %w", err)
	}
	return result, nil
}

// Search returns all labels whose names match the given query string.
func (s *LabelService) Search(ctx context.Context, name string) ([]Label, error) {
	results, err := s.repository.SearchByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("search labels: %w", err)
	}
	return results, nil
}
