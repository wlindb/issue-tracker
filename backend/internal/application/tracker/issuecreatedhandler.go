package tracker

import (
	"context"
	"log/slog"

	issuedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
)

// IssueCreatedHandler handles IssueCreatedEvent by logging its details.
type IssueCreatedHandler struct{}

// NewIssueCreatedHandler returns an IssueCreatedHandler.
func NewIssueCreatedHandler() *IssueCreatedHandler {
	return &IssueCreatedHandler{}
}

// Handle logs the details of the created issue.
func (h *IssueCreatedHandler) Handle(event issuedomain.IssueCreatedEvent) {
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

func (h *IssueCreatedHandler) Handler(ctx context.Context, event issuedomain.IssueCreatedEvent) {
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
