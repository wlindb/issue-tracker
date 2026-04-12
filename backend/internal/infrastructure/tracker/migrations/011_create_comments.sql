-- +goose Up
CREATE TABLE IF NOT EXISTS comments (
    id           UUID        PRIMARY KEY,
    body         TEXT        NOT NULL,
    author_id    UUID        NOT NULL,
    issue_id     UUID        NOT NULL REFERENCES issues(id),
    workspace_id UUID        NOT NULL REFERENCES workspaces(id),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE comments ENABLE ROW LEVEL SECURITY;
ALTER TABLE comments FORCE ROW LEVEL SECURITY;
CREATE POLICY workspace_isolation ON comments
  USING (
    workspace_id = NULLIF(current_setting('app.workspace_id', true), '')::uuid
    AND EXISTS (
      SELECT 1 FROM workspace_members
      WHERE workspace_id = NULLIF(current_setting('app.workspace_id', true), '')::uuid
        AND user_id      = NULLIF(current_setting('app.user_id', true), '')::uuid
    )
  );

GRANT SELECT, INSERT, UPDATE, DELETE ON comments TO appuser;

-- +goose Down
DROP POLICY IF EXISTS workspace_isolation ON comments;
ALTER TABLE comments DISABLE ROW LEVEL SECURITY;
DROP TABLE IF EXISTS comments;
