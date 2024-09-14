package balance

import (
	"context"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type AbstrStorage interface {
	SelectBalanceByLogin(ctx context.Context, login string) (int, int, error)
}

type BalanceHandler struct {
	DB AbstrStorage

	// сделано для того, чтобы мокать работу с токеном в тестах
	GetLogin func(context.Context) (jwt.Token, map[string]interface{}, error)
}

func InitHandler(db AbstrStorage) BalanceHandler {
	return BalanceHandler{
		DB:       db,
		GetLogin: jwtauth.FromContext,
	}
}
