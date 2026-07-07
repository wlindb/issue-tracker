-- +goose Up
CREATE EXTENSION IF NOT EXISTS pg_textsearch;

CREATE TABLE IF NOT EXISTS issue_documents (
    id           UUID        PRIMARY KEY,
    workspace_id UUID        NOT NULL,
    title        TEXT        NOT NULL,
    description  TEXT        NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS issue_documents_description_bm25_idx ON issue_documents
    USING bm25(description) WITH (text_config='english');

-- +goose Down
DROP INDEX IF EXISTS issue_documents_description_bm25_idx;
DROP TABLE IF EXISTS issue_documents;
DROP EXTENSION IF EXISTS pg_textsearch;
