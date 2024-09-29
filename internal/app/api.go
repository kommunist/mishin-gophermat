package app

import (
	"log/slog"
	"mishin-gophermat/internal/handlers/balance"
	"mishin-gophermat/internal/handlers/listorders"
	"mishin-gophermat/internal/handlers/listwithdrawns"
	"mishin-gophermat/internal/handlers/login"
	"mishin-gophermat/internal/handlers/postorders"
	"mishin-gophermat/internal/handlers/postwithdrawns"
	"mishin-gophermat/internal/handlers/registration"
	"mishin-gophermat/internal/middlewares/auth"
	"mishin-gophermat/internal/secure"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	jwtauth "github.com/go-chi/jwtauth/v5"
)

func (a *App) InitServer(h http.Handler) {
	a.srv = &http.Server{
		Addr:    a.Config.RunAddress,
		Handler: h,
	}
}

func (app *App) CollectHandlers() http.Handler {
	regH := registration.InitHandler(app.DB)
	poH := postorders.InitHandler(app.DB, app.Acr.AcrChan)
	loH := listorders.InitHandler(app.DB)
	balH := balance.InitHandler(app.DB)
	pwithdrawns := postwithdrawns.InitHandler(app.DB)
	lwithdrawns := listwithdrawns.InitHandler(app.DB)
	loginH := login.InitHandler(app.DB)
	authM := auth.InitAuth()

	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(secure.TokenAuth))
		r.Use(jwtauth.Authenticator(secure.TokenAuth))
		r.Use(authM.Auth)
		r.Post("/api/user/orders", poH.Process)
		r.Get("/api/user/orders", loH.Process)
		r.Get("/api/user/balance", balH.Process)
		r.Post("/api/user/balance/withdraw", pwithdrawns.Process)
		r.Get("/api/user/withdrawals", lwithdrawns.Process)
	})

	r.Post("/api/user/register", regH.Process)
	r.Post("/api/user/login", loginH.Process)

	return r
}

func (app *App) StartAPI() {
	app.InitServer(
		app.CollectHandlers(),
	)

	slog.Info("Start server on")

	err := app.srv.ListenAndServe()

	if err != nil {
		slog.Error("Server failed to start", "err", err)
		app.Finish(true)
	}

}
