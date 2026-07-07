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

type IssuePriorityUpdatedEvent struct {
	OccurredAt time.Time `json:"occurred_at"`
	Payload    Issue
}

var PriorityUpdated = event.New[IssuePriorityUpdatedEvent]()

type IssueTitleUpdatedEvent struct {
	OccurredAt time.Time `json:"occurred_at"`
	Payload    Issue
}

var TitleUpdated = event.New[IssueTitleUpdatedEvent]()

type IssueAssigneeUpdatedEvent struct {
	OccurredAt time.Time `json:"occurred_at"`
	Payload    Issue
}

var AssigneeUpdated = event.New[IssueAssigneeUpdatedEvent]()

type IssueDescriptionUpdatedEvent struct {
	OccurredAt time.Time `json:"occurred_at"`
	Payload    Issue
}

var DescriptionUpdated = event.New[IssueDescriptionUpdatedEvent]()

type IssueLabelAddedEvent struct {
	OccurredAt time.Time `json:"occurred_at"`
	Payload    Issue
}

var LabelAdded = event.New[IssueLabelAddedEvent]()
