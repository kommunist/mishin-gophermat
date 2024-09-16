package storage

import (
	"context"
	"log/slog"
)

func (db *DB) SelectBalanceByLogin(ctx context.Context, login string) (int, int, error) { // current, withdrawn
	var checkDebit interface{}
	var checkCredit interface{}
	var debit int
	var credit int

	row := db.driver.QueryRowContext(
		ctx,
		"SELECT sum(orders.value) FROM orders WHERE user_login = $1", login,
	)
	err := row.Scan(&checkDebit)
	if checkDebit != nil {
		debit = int(checkDebit.(int64))
	}
	if err != nil {
		slog.Info("error when scan data for debit", "err", err)
		return 0, 0, err
	}

	row = db.driver.QueryRowContext(
		ctx,
		"SELECT sum(withdrawns.value) FROM withdrawns WHERE user_login = $1", login,
	)
	err = row.Scan(&checkCredit)
	if checkCredit != nil {
		credit = int(checkCredit.(int64))
	}
	if err != nil {
		slog.Info("error when scan data for credit", "err", err)
		return 0, 0, err
	}

	return debit - credit, int(credit), nil
}
