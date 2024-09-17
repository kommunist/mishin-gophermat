package app

import "mishin-gophermat/internal/accrual"

func (app *App) InitAsync() {
	acr := accrual.InitAccrual(app.DB, app.Config.AccrualURI)
	acr.InitWorkers(app.AcrChan)
}
