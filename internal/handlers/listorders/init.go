package listorders

import (
	"context"
	"mishin-gophermat/internal/models"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type OrdersGetter interface {
	OrdersGet(ctx context.Context, login string) (data []models.Order, err error)
}

type ListOrdersHandler struct {
	DB OrdersGetter

	// сделано для того, чтобы мокать работу с токеном в тестах
	GetLogin func(context.Context) (jwt.Token, map[string]any, error)
}

func InitHandler(db OrdersGetter) ListOrdersHandler {
	return ListOrdersHandler{
		DB:       db,
		GetLogin: jwtauth.FromContext,
	}
}
