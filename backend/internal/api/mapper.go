package api

import (
	"github.com/wlindb/issue-tracker/internal/api/model"
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
	return trackerdomain.ListProjectQuery{
		Cursor: params.Cursor,
		Limit:  params.Limit,
	}
}
