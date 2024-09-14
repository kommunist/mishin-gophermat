-- +goose Up
-- +goose StatementBegin
CREATE TABLE withdrawns(
  id SERIAL PRIMARY KEY,
  number TEXT, -- кажется, что нигде не нужен
  user_login TEXT,
  uploaded_at TIMESTAMP DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE withdrawns;
-- +goose StatementEnd
