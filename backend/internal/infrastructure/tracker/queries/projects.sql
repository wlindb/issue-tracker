-- name: CreateProject :one
INSERT INTO projects (id, owner_id, name, description, created_at, updated_at)
VALUES (@id, @owner_id, @name, @description, NOW(), NOW())
RETURNING *;

-- name: ListProjects :many
SELECT * FROM projects
ORDER BY created_at DESC
LIMIT @project_limit::int4;
