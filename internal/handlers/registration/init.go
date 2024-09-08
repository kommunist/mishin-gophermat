package registration

import (
	"context"
)

type AbstrStorage interface {
	CreateUser(ctx context.Context, login string, password string) error
}

type RegistrationHandler struct {
	DB AbstrStorage
}

func InitHandler(db AbstrStorage) RegistrationHandler {
	return RegistrationHandler{DB: db}
}
