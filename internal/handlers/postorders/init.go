package postorders

import (
	"context"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type AbstrStorage interface {
	SelectOrderByNumber(ctx context.Context, numer string) (data map[string]interface{}, err error)
	CreateOrder(ctx context.Context, number string, userLogin string) error
}

type PostOrdersHandler struct {
	DB      AbstrStorage
	acrChan chan string

	// сделано для того, чтобы мокать работу с токеном в тестах
	GetLogin func(context.Context) (jwt.Token, map[string]interface{}, error)
}

func InitHandler(db AbstrStorage, acrChan chan string) PostOrdersHandler {
	return PostOrdersHandler{
		DB:      db,
		acrChan: acrChan,

		GetLogin: jwtauth.FromContext,
	}
}
