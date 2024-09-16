package storage

import (
	"context"
	"log/slog"
)

func (db *DB) CreateWithdrawn(ctx context.Context, userLogin string, number string, value int) error {
	tx, err := db.driver.Begin()

	if err != nil {
		slog.Error("Error when open transaction", "err", err)
		return err
	}

	data, err := tx.QueryContext(ctx,
		"INSERT INTO withdrawns (number, user_login) VALUES ($1, $2) RETURNING ID",
		number, userLogin,
	)

	if err != nil {
		slog.Error("Error when create withdrawn", "err", err)
		return err
	}

	var withdrawnId int
	err = data.Scan(&withdrawnId)

	if err != nil {
		slog.Error("Error when scan data", "err", err)
		return err
	}

	_, err = tx.ExecContext(
		ctx, "INSERT INTO balance_items (withdrawn_id, value) VALUES ($1, $2)", withdrawnId, value,
	)
	if err != nil {
		slog.Error("Error when insert balance_items", "err", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		slog.Error("Error when commit transaction", "err", err)
		return err
	}

	return nil
}

func (db *DB) SelectWithdrawnsByLogin(ctx context.Context, login string) ([]map[string]interface{}, error) {
	r := make([]map[string]interface{}, 0)
	var number, processedAt string
	var value int

	rows, err := db.driver.QueryContext(
		ctx,
		`
		SELECT number, processed_at, balance_items.value 
		FROM orders 
		LEFT JOIN balance_items on balance_items.withdrawn_id = withdrawns.id
		WHERE withdrawns.user_login = $1 limit 1`,
		login,
	)
	if err != nil {
		slog.Info("Error when get data from pg", "err", err)
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&number, &processedAt, &value)
		if err != nil {
			slog.Error("Error when scan data from result", "err", err)
			return nil, err
		}
		r = append(
			r,
			map[string]interface{}{
				"number":      number,
				"processedAt": processedAt,
				"value":       value,
			},
		)
	}
	return r, nil
}
