package issue

import (
	"time"

	"github.com/google/uuid"
)

// Status represents the workflow state of an issue.
type Status string

const (
	StatusBacklog    Status = "backlog"
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
	StatusCancelled  Status = "cancelled"
)

// Priority represents the importance level of an issue.
type Priority string

const (
	PriorityNone   Priority = "none"
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
	PriorityUrgent Priority = "urgent"
)

// Issue is the domain representation of a tracked issue.
type Issue struct {
	ID          uuid.UUID
	Identifier  string
	Title       string
	Description *string
	Status      Status
	Priority    Priority
	Labels      []string
	AssigneeID  *uuid.UUID
	ProjectID   uuid.UUID
	ReporterID  uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// IssuePage is the paginated result of a ListIssues operation.
type IssuePage struct {
	Items      []Issue
	NextCursor *string
}

// ListIssueQuery holds all parameters for listing issues within a project.
type ListIssueQuery struct {
	Cursor     *string
	Limit      *int
	Status     *Status
	Priority   *Priority
	AssigneeID *uuid.UUID
}

// CreateIssueRequest carries the input data for creating a new issue.
type CreateIssueRequest struct {
	Title       string
	Description *string
	Status      Status
	Priority    Priority
	AssigneeID  *uuid.UUID
}
