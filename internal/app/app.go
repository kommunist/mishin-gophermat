package app

import (
	"mishin-gophermat/internal/accrual"
	"mishin-gophermat/internal/config"
	"mishin-gophermat/internal/secure"
	"mishin-gophermat/internal/storage"
	"net/http"
)

type App struct {
	DB         *storage.DB       // база
	Config     config.MainConfig // конфиг
	srv        *http.Server      // api сервер
	Acr        *accrual.Accrual  // структура асинхронной работы с accrual
	FinishChan chan struct{}     // канал, что все завершили
}

func InitApp() *App {
	secure.InitSecure()

	c := config.MakeConfig()
	c.InitConfig()
	db := storage.MakeDB(c.DatabaseURI)

	acr := accrual.InitAccrual(&db, c.AccrualURI)

	a := App{
		Config:     c,
		DB:         &db,
		FinishChan: make(chan struct{}),
		Acr:        &acr,
	}
	WaitFinish(&a)

	return &a
}
