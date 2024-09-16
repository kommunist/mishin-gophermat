package listwithdrawns

import (
	"context"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type AbstrStorage interface {
	SelectWithdrawnsByLogin(ctx context.Context, login string) (data []map[string]interface{}, err error)
}

type ListWithdrawns struct {
	DB AbstrStorage

	// сделано для того, чтобы мокать работу с токеном в тестах
	GetLogin func(context.Context) (jwt.Token, map[string]interface{}, error)
}

func InitHandler(db AbstrStorage) ListWithdrawns {
	return ListWithdrawns{
		DB:       db,
		GetLogin: jwtauth.FromContext,
	}
}
