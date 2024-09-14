package storage

import (
	"context"
	"log/slog"
)

func (db *DB) SelectBalanceByLogin(ctx context.Context, login string) (int, int, error) { // current, withdrawn
	var current int
	var withdrawn int

	row := db.driver.QueryRowContext(
		ctx,
		`
		SELECT 
		SUM(CASE WHEN bi.withdrawn_id is not null THEN bi.value * -1 ELSE bi.value END) AS current,
		SUM(CASE WHEN bi.withdrawn_id is not null THEN bi.value ELSE 0 END) AS withdrawn,
		FROM balance_items bi
		INNER JOIN orders WHERE bi.order_id = orders.id
		INNER JOIN withdrawns where bi.withdrawn_id = withdrawns.id
		where orders.user_login = $1 or withdrawn.user_login = $1
		`, login,
	)
	err := row.Scan(&current, &withdrawn)

	if err != nil {
		slog.Info("error when scan data", "err", err)
		return 0, 0, err
	}

	return current, withdrawn, nil
}
