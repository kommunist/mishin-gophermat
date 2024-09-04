package storage

import (
	"context"
	"log/slog"
)

func (db *DB) CreateUser(ctx context.Context, login string, password string) error {
	_, err := db.driver.ExecContext(ctx,
		"INSERT INTO users (login, password) VALUES ($1, $2)",
		login, password,
	)

	if err != nil { // необходимо добавит обработку ситуации, когда пользователь уже есть
		slog.Error("Error when create user", "err", err)
		return err
	}

	return nil
}
