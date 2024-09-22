package auth

import (
	jwtauth "github.com/go-chi/jwtauth/v5"
)

var TokenAuth *jwtauth.JWTAuth

func InitAuth() {
	TokenAuth = jwtauth.New("HS256", []byte("secret"), nil) // пока просто слово secret
}

func Encrypt(data map[string]any) string {
	if TokenAuth == nil {
		InitAuth()
	}
	_, tokenString, _ := TokenAuth.Encode(data)

	return tokenString
}
