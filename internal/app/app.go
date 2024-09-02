package app

import (
	"mishin-gophermat/internal/config"
	"mishin-gophermat/internal/handlers/registration"
	"mishin-gophermat/internal/storage"
	"net/http"
)

type App struct {
	db     storage.DB
	config config.MainConfig
}

func MakeApp() App {
	c := config.MakeConfig()
	c.InitConfig()
	db := storage.MakeDB(c.DatabaseURI)
	return App{
		db:     db,
		config: c,
	}
}

func (h *App) Registration(w http.ResponseWriter, r *http.Request) {
	registration.Process(w, r, h.db)
}

func ErrorResponse(w *http.ResponseWriter) {
	wrt := *w
	http.Error(wrt, "Internal Server Error", http.StatusInternalServerError)
}
