package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPass(t *testing.T) {
	t.Run("when_correct_password", func(t *testing.T) {
		c := Crypt{}
		pass := "password"

		hashed, _ := c.PassHash(pass)
		checked := c.PassCheck(pass, hashed)

		assert.Equal(
			t,
			true,
			checked,
			"correct password must return true check",
		)

	})

	t.Run("when_incorrect_password", func(t *testing.T) {
		c := Crypt{}
		pass := "password"

		hashed, _ := c.PassHash(pass)
		checked := c.PassCheck("incorrect", hashed)

		assert.Equal(
			t,
			false,
			checked,
			"incorrect password must return false check",
		)

	})

}
