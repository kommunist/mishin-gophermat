package balance

import (
	"context"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type BalanceGetter interface {
	BalanceGet(ctx context.Context, login string) (float64, float64, error)
}

type BalanceHandler struct {
	DB BalanceGetter

	// сделано для того, чтобы мокать работу с токеном в тестах
	GetLogin func(context.Context) (jwt.Token, map[string]any, error)
}

func InitHandler(db BalanceGetter) BalanceHandler {
	return BalanceHandler{
		DB:       db,
		GetLogin: jwtauth.FromContext,
	}
}
