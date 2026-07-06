package embedding

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/wlindb/issue-tracker/internal/pkg/domain/event"
	"github.com/wlindb/issue-tracker/internal/pkg/tracker/model"
)

type Option func(*EmbeddingHandler) error

func WithIssueCreated(subscriber event.SubscriberOf[model.IssueCreatedEvent]) Option {
	return func(h *EmbeddingHandler) error {
		return subscriber.Subscribe(h.HandleIssueCreated)
	}
}

type EmbeddingHandler struct{}

func NewEmbeddingHandler(opts ...Option) (EmbeddingHandler, error) {
	var zero EmbeddingHandler
	handler := EmbeddingHandler{}
	for _, opt := range opts {
		if err := opt(&handler); err != nil {
			return zero, fmt.Errorf("embedding handler option: %w", err)
		}
	}
	return handler, nil
}

func (h EmbeddingHandler) HandleIssueCreated(_ context.Context, event model.IssueCreatedEvent) error {
	slog.Info("issue created",
		"issue_id", event.Payload.ID,
		"reporter_id", event.Payload.ReporterID,
		"title", event.Payload.Title,
		"status", event.Payload.Status,
		"priority", event.Payload.Priority,
		"occurred_at", event.OccurredAt,
	)

	return nil
}
