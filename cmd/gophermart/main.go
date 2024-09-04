package main

import (
	"log/slog"
	"mishin-gophermat/internal/app"
	"mishin-gophermat/internal/handlers/registration"
	"net/http"
	"os"

	chi "github.com/go-chi/chi/v5"
)

func main() {
	app := app.InitApp()

	regH := registration.InitHandler(app.Config, app.DB)

	r := chi.NewRouter()
	// r.Use(jwtauth.Verifier(tokenAuth))
	// r.Use(jwtauth.Authenticator(tokenAuth))

	r.Post("/api/user/register", regH.Process)

	slog.Info("Start server on")
	err := http.ListenAndServe(app.Config.RunAddress, r)
	if err != nil {
		slog.Error("Server failed to start", "err", err)
		os.Exit(1)
	}
}
