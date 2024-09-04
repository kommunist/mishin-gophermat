package app

import (
	"mishin-gophermat/internal/config"
	"mishin-gophermat/internal/storage"
)

type App struct {
	DB     *storage.DB
	Config config.MainConfig
}

func InitApp() App {
	c := config.MakeConfig()
	c.InitConfig()
	db := storage.MakeDB(c.DatabaseURI)

	return makeApp(c, &db)
}

func makeApp(c config.MainConfig, db *storage.DB) App {
	return App{Config: c, DB: db}
}
