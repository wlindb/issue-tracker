-- +goose Up
CREATE TABLE IF NOT EXISTS projects (
    id          UUID        PRIMARY KEY,
    owner_id    UUID        NOT NULL,
    name        TEXT        NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS projects;
