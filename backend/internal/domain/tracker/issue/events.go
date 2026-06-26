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

type IssueStatusUpdatedEvent struct {
	OccurredAt time.Time `json:"occurred_at"`
	Payload    Issue
}

var StatusUpdated = event.New[IssueStatusUpdatedEvent]()

type IssueTitleUpdatedEvent struct {
	OccurredAt time.Time `json:"occurred_at"`
	Payload    Issue
}

var TitleUpdated = event.New[IssueTitleUpdatedEvent]()
