package tracker

import (
	"github.com/wlindb/issue-tracker/internal/domain/tracker/comment"
	"github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
	"github.com/wlindb/issue-tracker/internal/domain/tracker/label"
	"github.com/wlindb/issue-tracker/internal/domain/tracker/project"
	"github.com/wlindb/issue-tracker/internal/pkg/tracker/model"
)

func ToLabelDTO(domain label.Label) model.Label {
	return model.Label{
		ID:   domain.ID,
		Name: domain.Name,
	}
}

func ToLabelsDTO(label []label.Label) []model.Label {
	dto := make([]model.Label, len(label))
	for i, domain := range label {
		dto[i] = ToLabelDTO(domain)
	}
	return dto
}

func ToIssueDTO(issue issue.Issue) model.Issue {
	return model.Issue{
		AssigneeID:  issue.AssigneeID,
		CreatedAt:   issue.CreatedAt,
		Description: issue.Description,
		ID:          issue.ID,
		Identifier:  issue.Identifier,
		Labels:      ToLabelsDTO(issue.Labels),
		Priority:    model.IssuePriority(issue.Priority),
		ProjectID:   issue.ProjectID,
		ReporterID:  issue.ReporterID,
		Status:      model.IssueStatus(issue.Status),
		Title:       issue.Title,
		UpdatedAt:   issue.UpdatedAt,
	}
}

func ToIssueCreatedEventDTO(event issue.IssueCreatedEvent) model.IssueCreatedEvent {
	return model.IssueCreatedEvent{
		OccurredAt: event.OccurredAt,
		Payload:    ToIssueDTO(event.Payload),
	}
}

func ToIssueStatusUpdatedEvent(event issue.IssueStatusUpdatedEvent) model.IssueStatusUpdatedEvent {
	return model.IssueStatusUpdatedEvent{
		OccurredAt: event.OccurredAt,
		Payload:    ToIssueDTO(event.Payload),
	}
}

func ToIssueTitleUpdatedEvent(event issue.IssueTitleUpdatedEvent) model.IssueTitleUpdatedEvent {
	return model.IssueTitleUpdatedEvent{
		OccurredAt: event.OccurredAt,
		Payload:    ToIssueDTO(event.Payload),
	}
}

func ToIssuePriorityUpdatedEvent(event issue.IssuePriorityUpdatedEvent) model.IssuePriorityUpdatedEvent {
	return model.IssuePriorityUpdatedEvent{
		OccurredAt: event.OccurredAt,
		Payload:    ToIssueDTO(event.Payload),
	}
}

func ToIssueAssigneeUpdatedEvent(event issue.IssueAssigneeUpdatedEvent) model.IssueAssigneeUpdatedEvent {
	return model.IssueAssigneeUpdatedEvent{
		OccurredAt: event.OccurredAt,
		Payload:    ToIssueDTO(event.Payload),
	}
}

func ToIssueDescriptionUpdatedEvent(event issue.IssueDescriptionUpdatedEvent) model.IssueDescriptionUpdatedEvent {
	return model.IssueDescriptionUpdatedEvent{
		OccurredAt: event.OccurredAt,
		Payload:    ToIssueDTO(event.Payload),
	}
}

func ToCommentDTO(comment comment.Comment) model.Comment {
	return model.Comment{
		ID:        comment.ID,
		Body:      comment.Body,
		AuthorID:  comment.AuthorID,
		IssueID:   comment.IssueID,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}
}

func ToCommentCreatedEvent(event comment.CommentCreatedEvent) model.CommentCreatedEvent {
	return model.CommentCreatedEvent{
		OccurredAt: event.OccurredAt,
		Payload:    ToCommentDTO(event.Payload),
	}
}

func ToProjectDTO(event project.Project) model.Project {
	return model.Project{
		ID:          event.ID,
		Identifier:  event.Identifier,
		Name:        event.Name,
		Description: event.Description,
		OwnerID:     event.OwnerID,
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	}
}

func ToProjectCreatedEvent(event project.ProjectCreatedEvent) model.ProjectCreatedEvent {
	return model.ProjectCreatedEvent{
		OccurredAt: event.OccurredAt,
		Payload:    ToProjectDTO(event.Payload),
	}
}
