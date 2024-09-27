package postwithdrawns

import (
	"context"
)

type WithdrawnCreator interface {
	WithdrawnCreate(ctx context.Context, userLogin string, number string, value float64) error
	BalanceGet(ctx context.Context, login string) (float64, float64, error)
}

type PostWithdrawsHandler struct {
	DB WithdrawnCreator
}

func InitHandler(db WithdrawnCreator) PostWithdrawsHandler {
	return PostWithdrawsHandler{
		DB: db,
	}
}
