package balance

import (
	"context"
)

type BalanceGetter interface {
	BalanceGet(ctx context.Context, login string) (float64, float64, error)
}

type BalanceHandler struct {
	DB BalanceGetter
}

func InitHandler(db BalanceGetter) BalanceHandler {
	return BalanceHandler{
		DB: db,
	}
}
