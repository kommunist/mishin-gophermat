package postwithdrawns

import (
	"context"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type WithdrawnCreator interface {
	WithdrawnCreate(ctx context.Context, userLogin string, number string, value float64) error
	BalanceGet(ctx context.Context, login string) (float64, float64, error)
}

type PostWithdrawsHandler struct {
	DB WithdrawnCreator

	// сделано для того, чтобы мокать работу с токеном в тестах
	GetLogin func(context.Context) (jwt.Token, map[string]interface{}, error)
}

func InitHandler(db WithdrawnCreator) PostWithdrawsHandler {
	return PostWithdrawsHandler{
		DB:       db,
		GetLogin: jwtauth.FromContext,
	}
}
