package listwithdrawns

import (
	"context"
	"mishin-gophermat/internal/models"
)

type WithdrawnsGetter interface {
	WithdrawnsGet(ctx context.Context, login string) (data []models.Withdrawn, err error)
}

type ListWithdrawns struct {
	DB WithdrawnsGetter
}

func InitHandler(db WithdrawnsGetter) ListWithdrawns {
	return ListWithdrawns{
		DB: db,
	}
}
