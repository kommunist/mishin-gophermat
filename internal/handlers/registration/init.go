package registration

import (
	"context"
	"mishin-gophermat/internal/config"
)

type abstrStorage interface {
	CreateUser(ctx context.Context, login string, password string) error
}

type RegistrationHandler struct {
	DB     abstrStorage
	Config config.MainConfig
}

func InitHandler(c config.MainConfig, db abstrStorage) RegistrationHandler {
	return RegistrationHandler{DB: db, Config: c}
}
