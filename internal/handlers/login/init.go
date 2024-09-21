package login

import (
	"context"
)

type UserGetter interface {
	UserGet(ctx context.Context, login string, password string) (bool, error) // login, error
}

type LoginHandler struct {
	DB UserGetter
}

func InitHandler(db UserGetter) LoginHandler {
	return LoginHandler{DB: db}
}
