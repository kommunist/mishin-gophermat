-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_orders_number on orders (number);
CREATE INDEX IF NOT EXISTS idx_orders_user_login on orders (user_login);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF NOT EXISTS idx_orders_number;
DROP INDEX IF NOT EXISTS idx_orders_user_login;
-- +goose StatementEnd
