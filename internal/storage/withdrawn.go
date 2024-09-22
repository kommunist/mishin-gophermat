package storage

import (
	"context"
	"log/slog"
)

func (db *DB) WithdrawnCreate(ctx context.Context, userLogin string, number string, value float64) error {
	_, err := db.driver.ExecContext(ctx,
		"INSERT INTO withdrawns (number, user_login, value) VALUES ($1, $2, $3)",
		number, userLogin, value,
	)

	if err != nil {
		slog.Error("Error when insert data", "err", err)
		return err
	}

	return nil
}

func (db *DB) WithdrawnsGet(ctx context.Context, login string) ([]map[string]any, error) {
	r := make([]map[string]any, 0)
	var number, processedAt string
	var value float64

	rows, err := db.driver.QueryContext(
		ctx,
		`
		SELECT number, processed_at, value FROM withdrawns WHERE withdrawns.user_login = $1 limit 1`,
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
			map[string]any{
				"number":      number,
				"processedAt": processedAt,
				"value":       value,
			},
		)
	}
	err = rows.Err()
	if err != nil {
		slog.Error("Error when iterate over rows", "err", err)
		return nil, err
	}

	return r, nil
}
