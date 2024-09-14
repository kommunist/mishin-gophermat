-- +goose Up
-- +goose StatementBegin
CREATE TABLE balance_items(
  id SERIAL PRIMARY KEY,
  value INTEGER,
  order_id INTEGER,
  withdrawn_id INTEGER
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE balance_items;
-- +goose StatementEnd
