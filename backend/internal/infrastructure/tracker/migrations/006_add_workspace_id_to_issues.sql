-- +goose Up
ALTER TABLE issues ADD COLUMN workspace_id UUID NOT NULL REFERENCES workspaces(id);
UPDATE issues i SET workspace_id = p.workspace_id FROM projects p WHERE i.project_id = p.id;
CREATE INDEX idx_issues_workspace_id ON issues(workspace_id);

-- +goose Down
DROP INDEX IF EXISTS idx_issues_workspace_id;
ALTER TABLE issues DROP COLUMN workspace_id;
