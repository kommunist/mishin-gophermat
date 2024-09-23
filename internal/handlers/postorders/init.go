package postorders

import (
	"context"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type OrderCreator interface {
	OrderByNumberGet(ctx context.Context, numer string) (login string, err error)
	OrderCreate(ctx context.Context, number string, userLogin string) error
}

type PostOrdersHandler struct {
	DB      OrderCreator
	acrChan chan string

	// сделано для того, чтобы мокать работу с токеном в тестах
	GetLogin func(context.Context) (jwt.Token, map[string]any, error)
}

func InitHandler(db OrderCreator, acrChan chan string) PostOrdersHandler {
	return PostOrdersHandler{
		DB:      db,
		acrChan: acrChan,

		GetLogin: jwtauth.FromContext,
	}
}
