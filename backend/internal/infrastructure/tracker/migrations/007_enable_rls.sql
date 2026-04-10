-- +goose Up
ALTER TABLE projects ENABLE ROW LEVEL SECURITY;
ALTER TABLE projects FORCE ROW LEVEL SECURITY;
CREATE POLICY workspace_isolation ON projects
  USING (
    workspace_id = NULLIF(current_setting('app.workspace_id', true), '')::uuid
    AND EXISTS (
      SELECT 1 FROM workspace_members
      WHERE workspace_id = NULLIF(current_setting('app.workspace_id', true), '')::uuid
        AND user_id      = NULLIF(current_setting('app.user_id', true), '')::uuid
    )
  );

ALTER TABLE issues ENABLE ROW LEVEL SECURITY;
ALTER TABLE issues FORCE ROW LEVEL SECURITY;
CREATE POLICY workspace_isolation ON issues
  USING (
    workspace_id = NULLIF(current_setting('app.workspace_id', true), '')::uuid
    AND EXISTS (
      SELECT 1 FROM workspace_members
      WHERE workspace_id = NULLIF(current_setting('app.workspace_id', true), '')::uuid
        AND user_id      = NULLIF(current_setting('app.user_id', true), '')::uuid
    )
  );

-- +goose Down
DROP POLICY workspace_isolation ON issues;
ALTER TABLE issues DISABLE ROW LEVEL SECURITY;

DROP POLICY workspace_isolation ON projects;
ALTER TABLE projects DISABLE ROW LEVEL SECURITY;
