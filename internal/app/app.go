package app

import (
	"mishin-gophermat/internal/config"
	"mishin-gophermat/internal/secure"
	"mishin-gophermat/internal/storage"
)

type App struct {
	DB      *storage.DB
	Config  config.MainConfig
	AcrChan chan string
}

func InitApp() App {
	secure.InitSecure()

	c := config.MakeConfig()
	c.InitConfig()
	db := storage.MakeDB(c.DatabaseURI)

	return App{
		Config:  c,
		DB:      &db,
		AcrChan: make(chan string, 5),
	}
}
