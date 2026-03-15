package project

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// StubRepository is a temporary no-op repository used until a real infrastructure
// implementation is wired in main.go.
type StubRepository struct{}

func (StubRepository) Create(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ string, _ *string) (*Project, error) {
	return nil, errors.New("not implemented")
}
