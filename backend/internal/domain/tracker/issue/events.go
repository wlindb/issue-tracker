package issue

import (
	"time"

	"github.com/google/uuid"
	"github.com/wlindb/issue-tracker/internal/pkg/domain/event"
)

// IssueCreatedEvent carries the data emitted when an issue is successfully persisted.
type IssueCreatedEvent struct {
	IssueID    uuid.UUID `json:"issue_id"`
	ProjectID  uuid.UUID `json:"project_id"`
	ReporterID uuid.UUID `json:"reporter_id"`
	Title      string    `json:"title"`
	Status     Status    `json:"status"`
	Priority   Priority  `json:"priority"`
	OccurredAt time.Time `json:"occurred_at"`
}

var Created = event.New[IssueCreatedEvent]()
