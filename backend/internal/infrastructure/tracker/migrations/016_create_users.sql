-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id    UUID PRIMARY KEY,
    email TEXT NOT NULL,
    name  TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS users;
