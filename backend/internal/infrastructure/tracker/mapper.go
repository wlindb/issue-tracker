package tracker

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	issuedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
	projectdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/project"
	trackerdb "github.com/wlindb/issue-tracker/internal/infrastructure/tracker/generated"
)

func rowToProject(row trackerdb.Project) *projectdomain.Project {
	p := &projectdomain.Project{
		ID:        row.ID,
		OwnerID:   row.OwnerID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
	if row.Description.Valid {
		s := row.Description.String
		p.Description = &s
	}
	return p
}

func rowsToProjects(rows []trackerdb.Project) []projectdomain.Project {
	projects := make([]projectdomain.Project, len(rows))
	for i, row := range rows {
		projects[i] = *rowToProject(row)
	}
	return projects
}

func createIssueParamsFromDomain(issue issuedomain.Issue) trackerdb.CreateIssueParams {
	var description pgtype.Text
	if issue.Description != nil {
		description = pgtype.Text{String: *issue.Description, Valid: true}
	}

	var assigneeID pgtype.UUID
	if issue.AssigneeID != nil {
		assigneeID = pgtype.UUID{Bytes: *issue.AssigneeID, Valid: true}
	}

	return trackerdb.CreateIssueParams{
		ID:          issue.ID,
		Identifier:  issue.Identifier,
		Title:       issue.Title,
		Description: description,
		Status:      string(issue.Status),
		Priority:    string(issue.Priority),
		Labels:      issue.Labels,
		AssigneeID:  assigneeID,
		ProjectID:   issue.ProjectID,
		ReporterID:  issue.ReporterID,
	}
}

func listIssuesParamsFromDomain(projectID uuid.UUID, query issuedomain.ListIssueQuery) trackerdb.ListIssuesParams {
	var status pgtype.Text
	if query.Status != nil {
		status = pgtype.Text{String: string(*query.Status), Valid: true}
	}

	var priority pgtype.Text
	if query.Priority != nil {
		priority = pgtype.Text{String: string(*query.Priority), Valid: true}
	}

	var assigneeID pgtype.UUID
	if query.AssigneeID != nil {
		assigneeID = pgtype.UUID{Bytes: *query.AssigneeID, Valid: true}
	}

	return trackerdb.ListIssuesParams{
		ProjectID:  projectID,
		Status:     status,
		Priority:   priority,
		AssigneeID: assigneeID,
	}
}

func issueToDomain(row trackerdb.Issue) *issuedomain.Issue {
	issue := &issuedomain.Issue{
		ID:         row.ID,
		Identifier: row.Identifier,
		Title:      row.Title,
		Status:     issuedomain.Status(row.Status),
		Priority:   issuedomain.Priority(row.Priority),
		Labels:     row.Labels,
		ProjectID:  row.ProjectID,
		ReporterID: row.ReporterID,
		CreatedAt:  row.CreatedAt.Time,
		UpdatedAt:  row.UpdatedAt.Time,
	}
	if row.Description.Valid {
		description := row.Description.String
		issue.Description = &description
	}
	if row.AssigneeID.Valid {
		assigneeID := uuid.UUID(row.AssigneeID.Bytes)
		issue.AssigneeID = &assigneeID
	}
	return issue
}

func issuesToDomain(rows []trackerdb.Issue) []issuedomain.Issue {
	issues := make([]issuedomain.Issue, len(rows))
	for idx, row := range rows {
		issues[idx] = *issueToDomain(row)
	}
	return issues
}
