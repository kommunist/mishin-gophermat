package main

import (
	"mishin-gophermat/internal/app"
)

func main() {
	app := app.InitApp()
	app.InitAsync()
	app.StartApi()
}
