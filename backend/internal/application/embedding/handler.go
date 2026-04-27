package embedding

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
	"github.com/wlindb/issue-tracker/internal/pkg/domain/event"
)

type IssueCreatedSubscriber interface {
	Subscribe(handler event.Subscriber[issue.IssueCreatedEvent]) error
}

type EmbeddingHandler struct{}

func NewEmbeddingHandler(issueCreatedSubscriber IssueCreatedSubscriber) (EmbeddingHandler, error) {
	var zero EmbeddingHandler

	handler := EmbeddingHandler{}
	if err := issueCreatedSubscriber.Subscribe(handler.Handler); err != nil {
		return zero, fmt.Errorf("subscribe issue created: %w", err)
	}

	return handler, nil
}

func (h EmbeddingHandler) Handler(_ context.Context, event issue.IssueCreatedEvent) {
	slog.Info("issue created",
		"issue_id", event.IssueID,
		"project_id", event.ProjectID,
		"reporter_id", event.ReporterID,
		"title", event.Title,
		"status", event.Status,
		"priority", event.Priority,
		"occurred_at", event.OccurredAt,
	)
}
