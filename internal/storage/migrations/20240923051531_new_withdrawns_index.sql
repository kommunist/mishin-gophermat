-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_withdrawns_user_login on withdrawns (user_login);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF NOT EXISTS idx_withdrawns_user_login;
-- +goose StatementEnd
