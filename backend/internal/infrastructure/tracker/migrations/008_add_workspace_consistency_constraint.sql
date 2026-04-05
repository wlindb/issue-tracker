-- +goose Up
-- Composite unique constraint required as the referencing target for the FK below.
ALTER TABLE projects ADD CONSTRAINT projects_workspace_id_id_unique UNIQUE (workspace_id, id);

-- Composite FK ensures an issue's workspace_id always matches its project's workspace_id.
-- FK checks bypass RLS, so this constraint holds even for direct DB writes.
ALTER TABLE issues ADD CONSTRAINT issues_workspace_matches_project
  FOREIGN KEY (workspace_id, project_id) REFERENCES projects (workspace_id, id);

-- +goose Down
ALTER TABLE issues DROP CONSTRAINT issues_workspace_matches_project;
ALTER TABLE projects DROP CONSTRAINT projects_workspace_id_id_unique;
