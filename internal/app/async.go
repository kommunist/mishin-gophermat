package app

func (app *App) InitAsync() {
	app.Acr.InitWorkers()
}
