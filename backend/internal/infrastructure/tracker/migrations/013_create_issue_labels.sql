-- +goose Up
CREATE TABLE IF NOT EXISTS issue_labels (
    issue_id UUID NOT NULL REFERENCES issues(id)  ON DELETE CASCADE,
    label_id UUID NOT NULL REFERENCES labels(id)  ON DELETE CASCADE,
    PRIMARY KEY (issue_id, label_id)
);

ALTER TABLE issue_labels ENABLE ROW LEVEL SECURITY;
ALTER TABLE issue_labels FORCE ROW LEVEL SECURITY;
CREATE POLICY workspace_isolation ON issue_labels
    USING (
        EXISTS (
            SELECT 1 FROM labels
            WHERE labels.id = issue_labels.label_id
              AND labels.workspace_id = NULLIF(current_setting('app.workspace_id', true), '')::uuid
        )
    );

GRANT SELECT, INSERT, DELETE ON issue_labels TO appuser;

CREATE OR REPLACE VIEW issue_with_labels
WITH (security_invoker = on)
AS
SELECT
    i.*,
    COALESCE(lbl.labels, '[]'::jsonb) AS labels
FROM issues i
LEFT JOIN LATERAL (
    SELECT jsonb_agg(l) AS labels
    FROM issue_labels il
    JOIN labels l
      ON l.id = il.label_id
    WHERE il.issue_id = i.id
) lbl on true;

-- +goose Down
DROP POLICY IF EXISTS workspace_isolation ON issue_labels;
ALTER TABLE issue_labels DISABLE ROW LEVEL SECURITY;
DROP TABLE IF EXISTS issue_labels;
DROP VIEW IF EXISTS issue_with_labels;
