package tracker

import (
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
