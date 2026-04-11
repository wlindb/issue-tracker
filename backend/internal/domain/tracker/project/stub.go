package project

import (
	"context"
	"errors"
)

// StubRepository is a temporary no-op repository used until a real infrastructure
// implementation is wired in main.go.
type StubRepository struct{}

func (StubRepository) Create(_ context.Context, _ Project) (Project, error) {
	return Project{}, errors.New("not implemented")
}

func (StubRepository) List(_ context.Context, _ ListProjectQuery) (Projects, error) {
	return Projects{}, errors.New("not implemented")
}
