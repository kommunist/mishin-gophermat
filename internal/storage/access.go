package storage

import (
	"context"
	"errors"
	"log/slog"
	"mishin-gophermat/internal/errors/exist"

	"github.com/lib/pq"
)

func (db *DB) CreateUser(ctx context.Context, login string, password string) error {
	_, err := db.driver.ExecContext(ctx,
		"INSERT INTO users (login, password) VALUES ($1, $2)",
		login, password,
	)

	var pqErr *pq.Error

	if errors.As(err, &pqErr) && pqErr.Code.Name() == "unique_violation" {
		slog.Info("User not uniq login", "login", login)
		return exist.NewExistError(err)
	}

	if err != nil { // необходимо добавит обработку ситуации, когда пользователь уже есть
		slog.Error("Error when create user", "err", err)
		return err
	}

	return nil
}
