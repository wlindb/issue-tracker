package api

import (
	"github.com/wlindb/issue-tracker/internal/api/model"
	trackerdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/project"
)

func projectToModel(p trackerdomain.Project) model.Project {
	return model.Project{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		OwnerId:     p.OwnerID,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func projectsToModel(projects []trackerdomain.Project) []model.Project {
	items := make([]model.Project, len(projects))
	for i, p := range projects {
		items[i] = projectToModel(p)
	}
	return items
}
