-- +goose Up
GRANT UPDATE ON labels TO appuser;

-- +goose Down
REVOKE UPDATE ON labels FROM appuser;
