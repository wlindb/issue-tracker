-- name: CreateProject :one
INSERT INTO projects (id, owner_id, name, description, workspace_id, created_at, updated_at)
VALUES (@id, @owner_id, @name, @description, current_setting('app.workspace_id')::uuid, NOW(), NOW())
RETURNING *;

-- name: ListProjects :many
SELECT * FROM projects
WHERE workspace_id = current_setting('app.workspace_id')::uuid
ORDER BY created_at DESC
LIMIT @project_limit::int4;

-- name: GetProject :one
SELECT * FROM projects
WHERE id = @id AND workspace_id = current_setting('app.workspace_id')::uuid;
