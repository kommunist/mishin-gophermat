package auth

import (
	"context"
	"mishin-gophermat/internal/secure"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type AuthStruct struct {
	// сделано для того, чтобы мокать работу с токеном в тестах
	GetLogin func(context.Context) (jwt.Token, map[string]any, error)
}

func InitAuth() AuthStruct {
	return AuthStruct{
		GetLogin: jwtauth.FromContext,
	}
}

func (a *AuthStruct) Auth(h http.Handler) http.Handler {
	authfn := func(w http.ResponseWriter, r *http.Request) {
		var currUser string

		_, claims, _ := a.GetLogin(r.Context())

		if userLogin := claims["login"]; userLogin != nil {
			currUser = claims["login"].(string)
		} else { // 401
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, secure.UserLoginKey, currUser)

		h.ServeHTTP(w, r.WithContext(ctx))

	}

	return http.HandlerFunc(authfn)

}
