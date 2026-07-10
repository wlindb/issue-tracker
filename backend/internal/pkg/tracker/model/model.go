package model

import (
	"time"

	"github.com/google/uuid"
)

type IssuePriority string

type IssueStatus string

type Label struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Issue struct {
	AssigneeID  *uuid.UUID    `json:"assigneeId,omitempty"`
	CreatedAt   time.Time     `json:"createdAt"`
	Description *string       `json:"description,omitempty"`
	ID          uuid.UUID     `json:"id"`
	Identifier  string        `json:"identifier"`
	Labels      []Label       `json:"labels"`
	Priority    IssuePriority `json:"priority"`
	ProjectID   uuid.UUID     `json:"projectId"`
	ReporterID  uuid.UUID     `json:"reporterId"`
	Status      IssueStatus   `json:"status"`
	Title       string        `json:"title"`
	UpdatedAt   time.Time     `json:"updatedAt"`
}

type IssueCreatedEvent struct {
	OccurredAt time.Time `json:"occurred_at"`
	Payload    Issue     `json:"payload"`
}

type IssueStatusUpdatedEvent struct {
	OccurredAt time.Time `json:"occurred_at"`
	Payload    Issue     `json:"payload"`
}

type IssueTitleUpdatedEvent struct {
	OccurredAt time.Time `json:"occurred_at"`
	Payload    Issue     `json:"payload"`
}

type IssuePriorityUpdatedEvent struct {
	OccurredAt time.Time `json:"occurred_at"`
	Payload    Issue     `json:"payload"`
}

type IssueAssigneeUpdatedEvent struct {
	OccurredAt time.Time `json:"occurred_at"`
	Payload    Issue     `json:"payload"`
}

type IssueDescriptionUpdatedEvent struct {
	OccurredAt time.Time `json:"occurred_at"`
	Payload    Issue     `json:"payload"`
}

type IssueLabelAddedEvent struct {
	OccurredAt time.Time `json:"occurred_at"`
	Payload    Issue     `json:"payload"`
}

type Comment struct {
	ID        uuid.UUID
	Body      string
	AuthorID  uuid.UUID
	IssueID   uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CommentCreatedEvent struct {
	OccurredAt time.Time `json:"occurred_at"`
	Payload    Comment   `json:"payload"`
}

type Project struct {
	ID          uuid.UUID
	Identifier  string
	Name        string
	Description *string
	OwnerID     uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ProjectCreatedEvent struct {
	OccurredAt time.Time `json:"occurred_at"`
	Payload    Project   `json:"payload"`
}
