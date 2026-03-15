-- name: CreateProject :one
INSERT INTO projects (id, owner_id, name, description, created_at, updated_at)
VALUES (@id, @owner_id, @name, @description, NOW(), NOW())
RETURNING *;
