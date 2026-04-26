package notification

import (
	"context"
	"log/slog"

	issuedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
)

type NotificationHandler struct{}

func NewNotificationHandler() NotificationHandler {
	return NotificationHandler{}
}

func (h NotificationHandler) Handler(_ context.Context, event issuedomain.IssueCreatedEvent) {
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
