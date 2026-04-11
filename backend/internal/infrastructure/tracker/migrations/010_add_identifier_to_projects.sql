-- +goose Up
ALTER TABLE projects ADD COLUMN identifier TEXT NOT NULL DEFAULT '';
UPDATE projects SET identifier = id::text WHERE identifier = '';
ALTER TABLE projects ALTER COLUMN identifier DROP DEFAULT;
ALTER TABLE projects ADD CONSTRAINT projects_workspace_id_identifier_key UNIQUE (workspace_id, identifier);

-- +goose Down
ALTER TABLE projects DROP CONSTRAINT IF EXISTS projects_workspace_id_identifier_key;
ALTER TABLE projects DROP COLUMN identifier;
