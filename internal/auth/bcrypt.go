package auth

import (
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type Crypt struct{}

func (c *Crypt) PassHash(pass string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Error when hash password", "err", err)
		return "", err
	}

	return string(hashed), nil
}

func (c *Crypt) PassCheck(pass string, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pass)) == nil
}
