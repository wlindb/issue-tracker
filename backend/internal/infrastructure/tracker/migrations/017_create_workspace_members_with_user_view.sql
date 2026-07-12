-- +goose Up
CREATE OR REPLACE VIEW workspace_members_with_user
WITH (security_invoker = on)
AS
SELECT wm.workspace_id, wm.user_id, wm.created_at, u.email, u.name
FROM workspace_members wm
JOIN users u ON u.id = wm.user_id;

-- +goose Down
DROP VIEW IF EXISTS workspace_members_with_user;
