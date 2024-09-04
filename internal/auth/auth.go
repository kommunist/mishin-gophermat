package auth

import (
	jwtauth "github.com/go-chi/jwtauth/v5"
)

var tokenAuth *jwtauth.JWTAuth

func InitAuth() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil) // пока просто слово secret
}

func Encrypt(data map[string]interface{}) string {
	if tokenAuth == nil {
		InitAuth()
	}
	_, tokenString, _ := tokenAuth.Encode(data)

	return tokenString
}
