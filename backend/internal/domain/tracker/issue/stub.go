package issue

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// StubRepository is a temporary no-op repository used until a real infrastructure
// implementation is wired in main.go.
type StubRepository struct{}

func (StubRepository) GetIssue(_ context.Context, _ uuid.UUID) (Issue, error) {
	return Issue{}, errors.New("not implemented")
}

func (StubRepository) ListIssues(_ context.Context, _ uuid.UUID, _ ListIssueQuery) (IssuePage, error) {
	return IssuePage{}, errors.New("not implemented")
}

func (StubRepository) CreateIssue(_ context.Context, _ Issue) (Issue, error) {
	return Issue{}, errors.New("not implemented")
}

func (StubRepository) Update(_ context.Context, _ Issue) (Issue, error) {
	return Issue{}, errors.New("not implemented")
}

// NoOpEventPublisher is a no-op EventPublisher used in tests and dev environments.
type NoOpEventPublisher struct{}

func (NoOpEventPublisher) PublishIssueCreated(_ IssueCreatedEvent) error {
	return nil
}
