package main

import (
	"mishin-gophermat/internal/app"
)

func main() {
	app := app.InitApp()
	app.InitAsync()
	app.StartAPI()

	<-app.FinishChan // ждем, когда процедура выключения закончится и закроет канал
}
