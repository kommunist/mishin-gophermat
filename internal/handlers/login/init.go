package login

import (
	"context"
	"mishin-gophermat/internal/secure"
)

type UserGetter interface {
	UserGet(ctx context.Context, login string) (string, error) // pass, error
}

type PassChecker interface {
	PassCheck(pass string, hashed string) bool
}

type LoginHandler struct {
	DB      UserGetter
	checker PassChecker
}

func InitHandler(db UserGetter) LoginHandler {
	return LoginHandler{
		DB:      db,
		checker: &secure.Crypt{},
	}
}
