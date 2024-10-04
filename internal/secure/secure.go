package secure

import (
	jwtauth "github.com/go-chi/jwtauth/v5"
)

type UserLoginKeyType int

const UserLoginKey UserLoginKeyType = 0

var TokenAuth *jwtauth.JWTAuth

func InitSecure() {
	TokenAuth = jwtauth.New("HS256", []byte("secret"), nil) // пока просто слово secret
}

func EncryptLogin(login string) string {
	if TokenAuth == nil {
		InitSecure()
	}

	_, tokenString, _ := TokenAuth.Encode(map[string]interface{}{"login": login})

	return tokenString
}
