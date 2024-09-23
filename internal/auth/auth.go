package auth

import (
	jwtauth "github.com/go-chi/jwtauth/v5"
)

var TokenAuth *jwtauth.JWTAuth

func InitAuth() {
	TokenAuth = jwtauth.New("HS256", []byte("secret"), nil) // пока просто слово secret
}

func EncryptLogin(login string) string {
	if TokenAuth == nil {
		InitAuth()
	}

	_, tokenString, _ := TokenAuth.Encode(map[string]interface{}{"login": login})

	return tokenString
}
