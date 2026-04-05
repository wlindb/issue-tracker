-- name: CreateIssue :one
INSERT INTO issues (id, identifier, title, description, status, priority, labels, assignee_id, project_id, reporter_id, workspace_id, created_at, updated_at)
VALUES (@id, @identifier, @title, @description, @status, @priority, @labels, @assignee_id, @project_id, @reporter_id, current_setting('app.workspace_id')::uuid, NOW(), NOW())
RETURNING *;

-- name: ListIssues :many
SELECT * FROM issues
WHERE project_id = @project_id
  AND workspace_id = current_setting('app.workspace_id')::uuid
  AND (sqlc.narg('status')::text      IS NULL OR status      = sqlc.narg('status'))
  AND (sqlc.narg('priority')::text    IS NULL OR priority    = sqlc.narg('priority'))
  AND (sqlc.narg('assignee_id')::uuid IS NULL OR assignee_id = sqlc.narg('assignee_id'))
ORDER BY created_at DESC
LIMIT 100;
