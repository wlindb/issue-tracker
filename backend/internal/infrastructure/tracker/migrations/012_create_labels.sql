-- +goose Up
CREATE TABLE IF NOT EXISTS labels (
    id           UUID        PRIMARY KEY,
    workspace_id UUID        NOT NULL REFERENCES workspaces(id),
    name         TEXT        NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (workspace_id, name)
);

CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX labels_name_trgm_idx ON labels USING GIN (name gin_trgm_ops);

ALTER TABLE issues DROP COLUMN labels;

ALTER TABLE labels ENABLE ROW LEVEL SECURITY;
ALTER TABLE labels FORCE ROW LEVEL SECURITY;
CREATE POLICY workspace_isolation ON labels
    USING (
        workspace_id = NULLIF(current_setting('app.workspace_id', true), '')::uuid
        AND EXISTS (
            SELECT 1 FROM workspace_members
            WHERE workspace_id = NULLIF(current_setting('app.workspace_id', true), '')::uuid
              AND user_id      = NULLIF(current_setting('app.user_id', true), '')::uuid
        )
    );

GRANT SELECT, INSERT ON labels TO appuser;

-- +goose Down
DROP POLICY IF EXISTS workspace_isolation ON labels;
ALTER TABLE labels DISABLE ROW LEVEL SECURITY;
DROP INDEX IF EXISTS labels_name_trgm_idx;
DROP TABLE IF EXISTS labels;
ALTER TABLE issues ADD COLUMN labels TEXT[] NOT NULL DEFAULT '{}';
