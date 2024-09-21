package storage

import (
	"context"
	"errors"
	"log/slog"
	"mishin-gophermat/internal/errors/exist"

	"github.com/lib/pq"
)

func (db *DB) UserCreate(ctx context.Context, login string, password string) error {
	_, err := db.driver.ExecContext(ctx,
		"INSERT INTO users (login, password) VALUES ($1, $2)",
		login, password,
	)

	var pqErr *pq.Error

	if errors.As(err, &pqErr) && pqErr.Code.Name() == "unique_violation" {
		slog.Info("User not uniq login", "login", login)
		return exist.NewExistError(err)
	}

	if err != nil {
		slog.Error("Error when create user", "err", err)
		return err
	}

	return nil
}

func (db *DB) UserGet(ctx context.Context, login string, password string) (bool, error) {
	row := db.driver.QueryRowContext(ctx,
		"select exists(select login from users where login = $1 and password = $2)", // наверное, не лучшая идея так логин/пароль впихивать
		login, password,
	)

	var ex bool

	err := row.Scan(&ex)
	if err != nil {
		slog.Error("Error when scan data", "err", err)
		return false, err
	}

	return ex, nil
}
