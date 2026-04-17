-- name: CreateComment :one
INSERT INTO comments (id, body, author_id, issue_id, workspace_id, created_at, updated_at)
VALUES (@id, @body, @author_id, @issue_id, current_setting('app.workspace_id')::uuid, NOW(), NOW())
RETURNING *;

-- name: ListCommentsByIssue :many
SELECT * FROM comments
WHERE issue_id = @issue_id
  AND workspace_id = current_setting('app.workspace_id')::uuid
ORDER BY created_at ASC;
