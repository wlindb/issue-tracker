package api

import (
	"github.com/wlindb/issue-tracker/internal/api/model"
	commentdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/comment"
	trackerdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/project"
)

func projectFromDomain(domain trackerdomain.Project) model.Project {
	return model.Project{
		Id:          domain.ID,
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
