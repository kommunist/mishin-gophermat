package listwithdrawns

import (
	"context"
	"mishin-gophermat/internal/models"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type WithdrawnsGetter interface {
	WithdrawnsGet(ctx context.Context, login string) (data []models.Withdrawn, err error)
}

type ListWithdrawns struct {
	DB WithdrawnsGetter

	// сделано для того, чтобы мокать работу с токеном в тестах
	GetLogin func(context.Context) (jwt.Token, map[string]any, error)
}

func InitHandler(db WithdrawnsGetter) ListWithdrawns {
	return ListWithdrawns{
		DB:       db,
		GetLogin: jwtauth.FromContext,
	}
}
