package registration

import (
	"context"
)

type UserCreator interface {
	UserCreate(ctx context.Context, login string, password string) error
}

type RegistrationHandler struct {
	DB UserCreator
}

func InitHandler(db UserCreator) RegistrationHandler {
	return RegistrationHandler{DB: db}
}
