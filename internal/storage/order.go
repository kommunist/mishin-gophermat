package storage

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"mishin-gophermat/internal/errors/exist"

	"github.com/lib/pq"
)

func (db *DB) SelectOrderByNumber(ctx context.Context, number string) (map[string]interface{}, error) {
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

func (db *DB) SelectOrdersByLogin(ctx context.Context, login string) ([]map[string]interface{}, error) {
	r := make([]map[string]interface{}, 0)
	var number, status, uploadedAt string
	var accrual interface{}

	rows, err := db.driver.QueryContext(
		ctx,
		`
		SELECT number, status, uploaded_at, balance_items.value 
		FROM orders 
		LEFT JOIN balance_items on balance_items.order_id = orders.id
		WHERE orders.user_login = $1 limit 1`,
		login,
	)
	if err != nil {
		slog.Info("Error when get data from pg", "err", err)
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&number, &status, &uploadedAt, &accrual)
		if err != nil {
			slog.Error("Error when scan data from result", "err", err)
			return nil, err
		}
		realAccrual := 0
		if accrual != nil {
			realAccrual = accrual.(int)
		}
		r = append(
			r,
			map[string]interface{}{
				"number":     number,
				"status":     status,
				"uploadedAt": uploadedAt,
				"accrual":    realAccrual,
			},
		)
	}
	return r, nil
}

func (db *DB) CreateOrder(ctx context.Context, number string, userLogin string) error {
	_, err := db.driver.ExecContext(ctx,
		"INSERT INTO orders (number, user_login) VALUES ($1, $2)",
		number, userLogin,
	)

	var pqErr *pq.Error

	if errors.As(err, &pqErr) && pqErr.Code.Name() == "unique_violation" {
		slog.Info("Order not uniq", "number", number)
		return exist.NewExistError(err)
	}

	if err != nil {
		slog.Error("Error when create order", "err", err)
		return err
	}

	return nil
}
