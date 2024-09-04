package main

import (
	"log/slog"
	"mishin-gophermat/internal/app"
	"mishin-gophermat/internal/auth"
	"net/http"
	"os"

	chi "github.com/go-chi/chi/v5"
	jwtauth "github.com/go-chi/jwtauth/v5"
)

var tokenAuth *jwtauth.JWTAuth

// func init() {
// 	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)

// 	// For debugging/example purposes, we generate and print
// 	// a sample jwt token with claims `user_id:123` here:
// 	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": 123})
// 	fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)
// }

func main() {

	app := app.MakeApp()
	auth.InitAuth()

	r := chi.NewRouter()
	// r.Use(jwtauth.Verifier(tokenAuth))
	// r.Use(jwtauth.Authenticator(tokenAuth))

	r.Post("/api/user/register", app.Registration)

	slog.Info("Start server on")
	err := http.ListenAndServe(app.Config.RunAddress, r)
	if err != nil {
		slog.Error("Server failed to start", "err", err)
		os.Exit(1)
	}
}
