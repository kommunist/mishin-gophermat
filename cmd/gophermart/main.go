package main

import (
	"log/slog"
	"mishin-gophermat/internal/app"
	"mishin-gophermat/internal/auth"
	"mishin-gophermat/internal/handlers/balance"
	"mishin-gophermat/internal/handlers/listorders"
	"mishin-gophermat/internal/handlers/listwithdrawns"
	"mishin-gophermat/internal/handlers/login"
	"mishin-gophermat/internal/handlers/postorders"
	"mishin-gophermat/internal/handlers/postwithdrawns"
	"mishin-gophermat/internal/handlers/registration"
	"net/http"
	"os"

	chi "github.com/go-chi/chi/v5"
	jwtauth "github.com/go-chi/jwtauth/v5"
)

func main() {
	app := app.InitApp()
	app.InitAsync()

	regH := registration.InitHandler(app.DB)
	poH := postorders.InitHandler(app.DB)
	loH := listorders.InitHandler(app.DB)
	balH := balance.InitHandler(app.DB)
	pwithdrawns := postwithdrawns.InitHandler(app.DB)
	lwithdrawns := listwithdrawns.InitHandler(app.DB)
	loginH := login.InitHandler(app.DB)

	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(auth.TokenAuth))
		r.Use(jwtauth.Authenticator(auth.TokenAuth))
		r.Post("/api/user/orders", poH.Process)
		r.Get("/api/user/orders", loH.Process)
		r.Get("/api/user/balance", balH.Process)
		r.Post("/api/user/balance/withdraw", pwithdrawns.Process)
		r.Get("/api/user/withdrawals", lwithdrawns.Process)
	})

	r.Post("/api/user/register", regH.Process)
	r.Post("/api/user/login", loginH.Process)

	slog.Info("Start server on")
	err := http.ListenAndServe(app.Config.RunAddress, r)
	if err != nil {
		slog.Error("Server failed to start", "err", err)
		os.Exit(1)
	}
}
