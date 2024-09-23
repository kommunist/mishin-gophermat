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

func (db *DB) UserGet(ctx context.Context, login string) (string, error) { // pass, err
	row := db.driver.QueryRowContext(ctx, "select password where login = $1", login)

	var pass string

	err := row.Scan(&pass)
	if err != nil {
		slog.Error("Error when scan data", "err", err)
		return "", err
	}

	return pass, nil
}
