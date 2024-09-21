package listwithdrawns

import (
	"context"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type WithdrawnsGetter interface {
	WithdrawnsGet(ctx context.Context, login string) (data []map[string]interface{}, err error)
}

type ListWithdrawns struct {
	DB WithdrawnsGetter

	// сделано для того, чтобы мокать работу с токеном в тестах
	GetLogin func(context.Context) (jwt.Token, map[string]interface{}, error)
}

func InitHandler(db WithdrawnsGetter) ListWithdrawns {
	return ListWithdrawns{
		DB:       db,
		GetLogin: jwtauth.FromContext,
	}
}
