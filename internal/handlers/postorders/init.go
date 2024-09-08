package postorders

import (
	"context"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type AbstrStorage interface {
	SelectOrder(ctx context.Context, numer string) (data map[string]interface{}, err error)
	CreateOrder(ctx context.Context, number string, userLogin string) error
}

type OrdersHandler struct {
	DB AbstrStorage

	// сделано для того, чтобы мокать работу с токеном в тестах
	GetLogin func(context.Context) (jwt.Token, map[string]interface{}, error)
}

func InitHandler(db AbstrStorage) OrdersHandler {
	return OrdersHandler{
		DB:       db,
		GetLogin: jwtauth.FromContext,
	}
}
