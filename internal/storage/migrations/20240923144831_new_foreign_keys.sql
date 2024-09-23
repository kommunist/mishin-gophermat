-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders ADD CONSTRAINT fk_orders_users FOREIGN KEY (user_login) REFERENCES users (login);
ALTER TABLE withdrawns ADD CONSTRAINT fk_orders_withdrawns FOREIGN KEY (user_login) REFERENCES users (login);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders DROP CONSTRAINT fk_orders_users;
ALTER TABLE withdrawns DROP CONSTRAINT fk_orders_withdrawns;
-- +goose StatementEnd
