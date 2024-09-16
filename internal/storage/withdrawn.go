package storage

import (
	"context"
	"log/slog"
)

func (db *DB) CreateWithdrawn(ctx context.Context, userLogin string, number string, value int) error {
	_, err := db.driver.ExecContext(ctx,
		"INSERT INTO withdrawns (number, user_login, value) VALUES ($1, $2, $3) RETURNING ID",
		number, userLogin, value,
	)

	if err != nil {
		slog.Error("Error when insert data", "err", err)
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
		SELECT number, processed_at, value FROM orders WHERE withdrawns.user_login = $1 limit 1`,
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
