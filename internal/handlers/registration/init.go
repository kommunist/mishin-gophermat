package registration

import (
	"context"
	"mishin-gophermat/internal/secure"
)

type UserCreator interface {
	UserCreate(ctx context.Context, login string, password string) error
}

type PassHasher interface {
	PassHash(pass string) (string, error)
}

type RegistrationHandler struct {
	DB     UserCreator
	hasher PassHasher
}

func InitHandler(db UserCreator) RegistrationHandler {
	return RegistrationHandler{
		DB:     db,
		hasher: &secure.Crypt{},
	}
}
