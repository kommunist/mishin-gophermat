package main

import (
	"log/slog"
	"mishin-gophermat/internal/app"
	"mishin-gophermat/internal/auth"
	"mishin-gophermat/internal/handlers/postorders"
	"mishin-gophermat/internal/handlers/registration"
	"net/http"
	"os"

	chi "github.com/go-chi/chi/v5"
	jwtauth "github.com/go-chi/jwtauth/v5"
)

func main() {
	app := app.InitApp()
	auth.InitAuth()

	regH := registration.InitHandler(app.DB)
	orderH := postorders.InitHandler(app.DB)

	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(auth.TokenAuth))
		r.Use(jwtauth.Authenticator(auth.TokenAuth))
		r.Post("/api/user/orders", orderH.Process)
	})

	r.Post("/api/user/register", regH.Process)

	slog.Info("Start server on")
	err := http.ListenAndServe(app.Config.RunAddress, r)
	if err != nil {
		slog.Error("Server failed to start", "err", err)
		os.Exit(1)
	}
}
