package main

import (
	"log/slog"
	"mishin-gophermat/internal/app"
	"net/http"
	"os"

	chi "github.com/go-chi/chi/v5"
)

func main() {

	app := app.MakeApp()

	r := chi.NewRouter()
	r.Post("/", app.Registration)

	slog.Info("Start server on")
	err := http.ListenAndServe("0.0.0.0:8080", r)
	if err != nil {
		slog.Error("Server failed to start", "err", err)
		os.Exit(1)
	}
}
