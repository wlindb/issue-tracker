package issue

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	nonSlugCharPattern    = regexp.MustCompile(`[^a-z0-9-]`)
	multipleDashesPattern = regexp.MustCompile(`-{2,}`)
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

type CreateIssueCommand struct {
	ProjectID   uuid.UUID
	ReporterID  uuid.UUID
	Title       string
	Description *string
	Status      Status
	Priority    Priority
	AssigneeID  *uuid.UUID
}

func (c CreateIssueCommand) Slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = nonSlugCharPattern.ReplaceAllString(s, "")
	s = multipleDashesPattern.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

type Slugifier func(s string) string

func (c CreateIssueCommand) ToIssue(id uuid.UUID, slugifier Slugifier) Issue {
	return Issue{
		ID:          id,
		Identifier:  slugifier(fmt.Sprintf("%s-%s", c.Title, id.String()[:8])),
		ProjectID:   c.ProjectID,
		ReporterID:  c.ReporterID,
		Title:       c.Title,
		Description: c.Description,
		Status:      c.Status,
		Priority:    c.Priority,
		Labels:      []string{},
		AssigneeID:  c.AssigneeID,
	}
}

type IssueRepository interface {
	ListIssues(ctx context.Context, projectID uuid.UUID, query ListIssueQuery) (IssuePage, error)
	CreateIssue(ctx context.Context, issue Issue) (*Issue, error)
}

var (
	ErrIssueNotFound   = errors.New("issue not found")
	ErrInvalidIssue    = errors.New("invalid issue")
	ErrNotImplemented  = errors.New("not implemented")
)
