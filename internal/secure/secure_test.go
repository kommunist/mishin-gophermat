package secure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	t.Run("encrypt_data", func(t *testing.T) {
		InitSecure()

		encrypted := EncryptLogin("marks")

		assert.Equal(
			t,
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpbiI6Im1hcmtzIn0.YbmB9Pl2FPkQ-8Xjo4GmGHYw9PjzO7LKK_8JATmtVrU",
			encrypted,
			"data must be correct encrypted",
		)

	})
}
