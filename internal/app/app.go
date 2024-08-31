package app

import (
	"mishin-gophermat/internal/handlers/registration"
	"net/http"
)

type App struct {
	db     string
	config string
}

func MakeApp() App {
	return App{
		db:     "here will be db",
		config: "here will be config",
	}
}

func (h *App) Registration(w http.ResponseWriter, r *http.Request) {
	registration.Process(w, r, h.db)
}

func ErrorResponse(w *http.ResponseWriter) {
	wrt := *w
	http.Error(wrt, "Internal Server Error", http.StatusInternalServerError)
}
