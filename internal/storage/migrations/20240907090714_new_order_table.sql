-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders(
  id SERIAL PRIMARY KEY,
  number TEXT UNIQUE,
  status TEXT default 'NEW',
  user_login TEXT,
  uploaded_at TIMESTAMP DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;
-- +goose StatementEnd
