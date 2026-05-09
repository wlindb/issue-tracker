package issue

import (
	"time"

	"github.com/wlindb/issue-tracker/internal/pkg/domain/event"
)

type IssueCreatedEvent struct {
	OccurredAt time.Time `json:"occurred_at"`
	Payload    Issue
}

var Created = event.New[IssueCreatedEvent]()
