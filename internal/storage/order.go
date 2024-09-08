package storage

import (
	"context"
	"errors"
	"log/slog"
	"mishin-gophermat/internal/errors/exist"

	"github.com/lib/pq"
)

<<<<<<< Updated upstream:internal/storage/access.go
func (db *DB) CreateUser(ctx context.Context, login string, password string) error {
	_, err := db.driver.ExecContext(ctx,
		"INSERT INTO users (login, password) VALUES ($1, $2)",
		login, password,
=======
func (db *DB) SelectOrder(ctx context.Context, number string) (map[string]interface{}, error) {
	var userLogin string

	r := make(map[string]interface{})
	row := db.driver.QueryRowContext(ctx, "SELECT user_login FROM orders where number = $1 limit 1", number)
	err := row.Scan(&userLogin)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		slog.Info("error when scan data", "err", err)
		return nil, err
	}

	r["userLogin"] = userLogin

	return r, nil
}

func (db *DB) CreateOrder(ctx context.Context, number string, userLogin string) error {
	_, err := db.driver.ExecContext(ctx,
		"INSERT INTO orders (number, user_login) VALUES ($1, $2)",
		number, userLogin,
>>>>>>> Stashed changes:internal/storage/order.go
	)

	var pqErr *pq.Error

	if errors.As(err, &pqErr) && pqErr.Code.Name() == "unique_violation" {
<<<<<<< Updated upstream:internal/storage/access.go
		slog.Info("User not uniq login", "login", login)
		return exist.NewExistError(err)
	}

	if err != nil { // необходимо добавит обработку ситуации, когда пользователь уже есть
		slog.Error("Error when create user", "err", err)
=======
		slog.Info("Order not uniq", "number", number)
		return exist.NewExistError(err)
	}

	if err != nil {
		slog.Error("Error when create order", "err", err)
>>>>>>> Stashed changes:internal/storage/order.go
		return err
	}

	return nil
}
