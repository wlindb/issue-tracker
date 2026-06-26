-- name: InsertIssueLabel :exec
INSERT INTO issue_labels (issue_id, label_id) VALUES (@issue_id, @label_id)
RETURNING *;

-- name: GetLabel :one
SELECT * FROM labels
WHERE id = @id
  AND workspace_id = current_setting('app.workspace_id')::uuid;
