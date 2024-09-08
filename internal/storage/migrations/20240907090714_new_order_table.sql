-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders(
  number TEXT PRIMARY KEY,
  status INTEGER default 0,
  accrual INTEGER default 0,
  user_login TEXT,
  updated_at TIMESTAMP DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;
-- +goose StatementEnd
