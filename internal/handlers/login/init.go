package login

import (
	"context"
)

type AbstrStorage interface {
	SelectUser(ctx context.Context, login string, password string) (bool, error) // login, error
}

type LoginHandler struct {
	DB AbstrStorage
}

func InitHandler(db AbstrStorage) LoginHandler {
	return LoginHandler{DB: db}
}
