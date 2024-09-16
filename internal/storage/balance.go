package storage

import (
	"context"
	"log/slog"
)

func (db *DB) SelectBalanceByLogin(ctx context.Context, login string) (int, int, error) { // current, withdrawn
	var current interface{}
	realCurrent := 0
	var withdrawn interface{}
	realWithdrawn := 0

	row := db.driver.QueryRowContext(
		ctx,
		`
		SELECT 
		SUM(CASE WHEN bi.withdrawn_id is not null THEN bi.value * -1 ELSE bi.value END) AS current,
		SUM(CASE WHEN bi.withdrawn_id is not null THEN bi.value ELSE 0 END) AS withdrawn
		FROM balance_items bi
		INNER JOIN orders ON bi.order_id = orders.id
		INNER JOIN withdrawns ON bi.withdrawn_id = withdrawns.id
		WHERE orders.user_login = $1 or withdrawns.user_login = $1
		`, login,
	)
	err := row.Scan(&current, &withdrawn)
	if current != nil {
		realCurrent = current.(int)
	}
	if withdrawn != nil {
		realWithdrawn = withdrawn.(int)
	}

	if err != nil {
		slog.Info("error when scan data", "err", err)
		return 0, 0, err
	}

	return realCurrent, realWithdrawn, nil
}
