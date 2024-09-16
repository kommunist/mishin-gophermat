package postwithdrawns

import (
	"context"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type AbstrStorage interface {
	CreateWithdrawn(ctx context.Context, userLogin string, number string, value float64) error
	SelectBalanceByLogin(ctx context.Context, login string) (float64, float64, error)
}

type PostWithdrawsHandler struct {
	DB AbstrStorage

	// сделано для того, чтобы мокать работу с токеном в тестах
	GetLogin func(context.Context) (jwt.Token, map[string]interface{}, error)
}

func InitHandler(db AbstrStorage) PostWithdrawsHandler {
	return PostWithdrawsHandler{
		DB:       db,
		GetLogin: jwtauth.FromContext,
	}
}
