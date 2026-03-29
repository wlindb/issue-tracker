-- +goose Up
CREATE TABLE IF NOT EXISTS issues (
    id          UUID        PRIMARY KEY,
    identifier  TEXT        NOT NULL,
    title       TEXT        NOT NULL,
    description TEXT,
    status      TEXT        NOT NULL,
    priority    TEXT        NOT NULL,
    labels      TEXT[]      NOT NULL,
    assignee_id UUID,
    project_id  UUID        NOT NULL REFERENCES projects(id),
    reporter_id UUID        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (project_id, identifier)
);

-- +goose Down
DROP TABLE IF EXISTS issues;
