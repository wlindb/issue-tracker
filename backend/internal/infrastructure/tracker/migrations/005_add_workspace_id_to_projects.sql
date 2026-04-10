-- +goose Up
ALTER TABLE projects ADD COLUMN workspace_id UUID NOT NULL REFERENCES workspaces(id);
CREATE INDEX idx_projects_workspace_id ON projects(workspace_id);

-- +goose Down
DROP INDEX IF EXISTS idx_projects_workspace_id;
ALTER TABLE projects DROP COLUMN workspace_id;
