-- +goose Up
-- +goose StatementBegin
CREATE TABLE users(
  user_id SERIAL PRIMARY KEY,
  login TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
