package app

import (
	"mishin-gophermat/internal/accrual"
	"mishin-gophermat/internal/auth"
	"mishin-gophermat/internal/config"
	"mishin-gophermat/internal/storage"
)

type App struct {
	DB      *storage.DB
	Config  config.MainConfig
	AcrChan chan string
}

func InitApp() App {
	auth.InitAuth()

	c := config.MakeConfig()
	c.InitConfig()
	db := storage.MakeDB(c.DatabaseURI)

	return App{
		Config:  c,
		DB:      &db,
		AcrChan: make(chan string, 5),
	}
}

func (app *App) InitAsync() {
	acr := accrual.InitAccrual(app.DB, app.Config.AccrualURI)
	acr.InitWorkers(app.AcrChan)
}
