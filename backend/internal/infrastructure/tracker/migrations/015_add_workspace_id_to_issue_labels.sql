-- +goose Up
ALTER TABLE issues ADD CONSTRAINT issues_workspace_id_id_unique UNIQUE (workspace_id, id);
ALTER TABLE labels ADD CONSTRAINT labels_workspace_id_id_unique UNIQUE (workspace_id, id);

ALTER TABLE issue_labels ADD COLUMN workspace_id UUID;
UPDATE issue_labels il SET workspace_id = i.workspace_id FROM issues i WHERE il.issue_id = i.id;
ALTER TABLE issue_labels ALTER COLUMN workspace_id SET NOT NULL;

-- Composite FKs ensure an issue_labels row's workspace_id always matches both
-- the referenced issue's and the referenced label's workspace_id. FK checks
-- bypass RLS, so these constraints hold even for direct DB writes.
ALTER TABLE issue_labels
  ADD CONSTRAINT issue_labels_workspace_matches_issue
    FOREIGN KEY (workspace_id, issue_id) REFERENCES issues (workspace_id, id) ON DELETE CASCADE,
  ADD CONSTRAINT issue_labels_workspace_matches_label
    FOREIGN KEY (workspace_id, label_id) REFERENCES labels (workspace_id, id) ON DELETE CASCADE;

DROP POLICY IF EXISTS workspace_isolation ON issue_labels;
CREATE POLICY workspace_isolation ON issue_labels
  USING (
    workspace_id = NULLIF(current_setting('app.workspace_id', true), '')::uuid
    AND EXISTS (
      SELECT 1 FROM workspace_members
      WHERE workspace_id = NULLIF(current_setting('app.workspace_id', true), '')::uuid
        AND user_id      = NULLIF(current_setting('app.user_id', true), '')::uuid
    )
  );

-- +goose Down
DROP POLICY IF EXISTS workspace_isolation ON issue_labels;
CREATE POLICY workspace_isolation ON issue_labels
USING (
    EXISTS (
        SELECT 1 FROM labels
        WHERE labels.id = issue_labels.label_id
            AND labels.workspace_id = NULLIF(current_setting('app.workspace_id', true), '')::uuid
    )
);

ALTER TABLE issue_labels DROP CONSTRAINT issue_labels_workspace_matches_label;
ALTER TABLE issue_labels DROP CONSTRAINT issue_labels_workspace_matches_issue;
ALTER TABLE issue_labels DROP COLUMN workspace_id;

ALTER TABLE labels DROP CONSTRAINT labels_workspace_id_id_unique;
ALTER TABLE issues DROP CONSTRAINT issues_workspace_id_id_unique;
