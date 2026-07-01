-- name: InsertIssueLabel :exec
INSERT INTO issue_labels (issue_id, label_id) VALUES (@issue_id, @label_id)
RETURNING *;

-- name: GetLabel :one
SELECT * FROM labels
WHERE id = @id
  AND workspace_id = current_setting('app.workspace_id')::uuid;

-- name: UpsertLabel :one
INSERT INTO labels (id, workspace_id, name, created_at)
VALUES (@id, current_setting('app.workspace_id')::uuid, @name, NOW())
ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name
RETURNING *;

-- name: ListLabelsByIDs :many
SELECT * FROM labels
WHERE id = ANY(@ids::uuid[]);

-- name: SearchLabelsByName :many
SELECT * FROM labels
WHERE name ILIKE '%' || @search || '%'
ORDER BY name;
