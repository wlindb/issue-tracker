package api

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/wlindb/issue-tracker/internal/api/model"
	commentdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/comment"
	issuedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
	trackerdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/project"
	workspacedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/workspace"
)

func projectFromDomain(domain trackerdomain.Project) model.Project {
	return model.Project{
		Id:          domain.ID,
		Identifier:  domain.Identifier,
		Name:        domain.Name,
		Description: domain.Description,
		OwnerId:     domain.OwnerID,
		CreatedAt:   domain.CreatedAt,
		UpdatedAt:   domain.UpdatedAt,
	}
}

func projectsFromDomain(domain []trackerdomain.Project) []model.Project {
	items := make([]model.Project, len(domain))
	for i, p := range domain {
		items[i] = projectFromDomain(p)
	}
	return items
}

func workspaceToDomain(ownerID uuid.UUID, req model.CreateWorkspaceRequest) (workspacedomain.Workspace, error) {
	workspace, err := workspacedomain.New(uuid.New(), req.Name, ownerID)
	if err != nil {
		return workspacedomain.Workspace{}, fmt.Errorf("workspace: %w", err)
	}
	return workspace, nil
}

func workspaceFromDomain(domain workspacedomain.Workspace) model.Workspace {
	return model.Workspace{
		Id:        domain.ID,
		Name:      domain.Name,
		OwnerId:   domain.OwnerID,
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
	}
}

func workspacesFromDomain(domain []workspacedomain.Workspace) []model.Workspace {
	items := make([]model.Workspace, len(domain))
	for i, w := range domain {
		items[i] = workspaceFromDomain(w)
	}
	return items
}

func workspaceMembersFromDomain(domain workspacedomain.WorkspaceMembers) []model.WorkspaceMember {
	items := make([]model.WorkspaceMember, len(domain.Members))
	for i, m := range domain.Members {
		items[i] = model.WorkspaceMember{
			Id: m.UserID,
		}
	}
	return items
}

func listProjectQueryFromRequest(params model.ListProjectsParams) trackerdomain.ListProjectQuery {
	return trackerdomain.NewListProjectQuery(params.Cursor, params.Limit)
}

func commentFromDomain(c commentdomain.Comment) model.Comment {
	return model.Comment{
		Id:        c.ID,
		Body:      c.Body,
		AuthorId:  c.AuthorID,
		IssueId:   c.IssueID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func commentsFromDomain(items []commentdomain.Comment) []model.Comment {
	out := make([]model.Comment, len(items))
	for i, c := range items {
		out[i] = commentFromDomain(c)
	}
	return out
}

func listCommentQueryFromRequest(params model.ListCommentsParams) commentdomain.ListCommentQuery {
	return commentdomain.NewListCommentQuery(params.Cursor, params.Limit)
}

func issueFromDomain(d issuedomain.Issue) model.Issue {
	labels := d.Labels
	if labels == nil {
		labels = []string{}
	}
	return model.Issue{
		Id:          d.ID,
		Identifier:  d.Identifier,
		Title:       d.Title,
		Description: d.Description,
		Status:      model.IssueStatus(d.Status),
		Priority:    model.IssuePriority(d.Priority),
		Labels:      labels,
		AssigneeId:  d.AssigneeID,
		ProjectId:   d.ProjectID,
		ReporterId:  d.ReporterID,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

func issuesFromDomain(domain []issuedomain.Issue) []model.Issue {
	items := make([]model.Issue, len(domain))
	for i, issue := range domain {
		items[i] = issueFromDomain(issue)
	}
	return items
}

func listIssueQueryFromRequest(params model.ListIssuesParams) issuedomain.ListIssueQuery {
	var status *issuedomain.Status
	if params.Status != nil {
		s := issuedomain.Status(*params.Status)
		status = &s
	}
	var priority *issuedomain.Priority
	if params.Priority != nil {
		p := issuedomain.Priority(*params.Priority)
		priority = &p
	}
	return issuedomain.ListIssueQuery{
		Cursor:     params.Cursor,
		Limit:      params.Limit,
		Status:     status,
		Priority:   priority,
		AssigneeID: params.AssigneeId,
	}
}

func createIssueCommandFromModel(projectID uuid.UUID, reporterID uuid.UUID, req model.CreateIssueRequest) issuedomain.CreateIssueCommand {
	return issuedomain.CreateIssueCommand{
		ProjectID:   projectID,
		ReporterID:  reporterID,
		Title:       req.Title,
		Description: req.Description,
		Status:      issuedomain.Status(req.Status),
		Priority:    issuedomain.Priority(req.Priority),
		AssigneeID:  req.AssigneeId,
	}
}
