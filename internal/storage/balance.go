package storage

import (
	"context"
	"log/slog"
)

func (db *DB) BalanceGet(ctx context.Context, login string) (float64, float64, error) { // current, withdrawn
	var checkBalance any
	var checkCredit any
	var balance float64
	var credit float64

	row := db.driver.QueryRowContext(
		ctx,
		`
		select sum(balance), sum(withdrawns) from (
	    	select (value * -1) as balance, value as withdrawns from withdrawns where user_login = $1
    		union all
    		select value as balance, 0 as withdrawns from orders where user_login = $1
		)
		`, login,
	)

	err := row.Scan(&checkBalance, &checkCredit)

	if err != nil {
		slog.Info("error when scan data for balance", "err", err)
		return 0.0, 0.0, err
	}

	if checkBalance != nil {
		balance = checkBalance.(float64)
	}
	if checkCredit != nil {
		credit = checkCredit.(float64)
	}
	return balance, credit, nil
}
