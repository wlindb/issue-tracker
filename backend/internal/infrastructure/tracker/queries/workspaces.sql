-- name: CreateWorkspace :one
INSERT INTO workspaces (id, name, owner_id, created_at, updated_at)
VALUES (@id, @name, @owner_id, NOW(), NOW())
RETURNING *;

-- name: InsertWorkspaceMember :exec
INSERT INTO workspace_members (workspace_id, user_id)
VALUES (@workspace_id, @user_id);

-- name: GetWorkspace :one
SELECT * FROM workspaces WHERE id = @id;

-- name: IsMember :one
SELECT EXISTS(
    SELECT 1 FROM workspace_members
    WHERE workspace_id = @workspace_id AND user_id = @user_id
) AS is_member;

-- name: ListWorkspacesForUser :many
SELECT w.* FROM workspaces w
JOIN workspace_members m ON m.workspace_id = w.id
WHERE m.user_id = @user_id
ORDER BY w.created_at DESC;
