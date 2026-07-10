-- name: UpsertUser :one
INSERT INTO users (id, email, name)
VALUES (@id, @email, @name)
ON CONFLICT (id) DO UPDATE
SET email = EXCLUDED.email,
    name  = EXCLUDED.name
RETURNING *;
