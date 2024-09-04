package registration

import (
	"context"
	"mishin-gophermat/internal/config"
)

type AbstrStorage interface {
	CreateUser(ctx context.Context, login string, password string) error
}

type RegistrationHandler struct {
	DB     AbstrStorage
	Config config.MainConfig
}

func InitHandler(c config.MainConfig, db AbstrStorage) RegistrationHandler {
	return RegistrationHandler{DB: db, Config: c}
}
