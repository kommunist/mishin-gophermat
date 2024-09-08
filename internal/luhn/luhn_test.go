package luhn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValid(t *testing.T) {
	t.Run("valid_string", func(t *testing.T) {
		res, _ := Valid([]byte("98265820"))

		assert.Equal(t, true, res)
	})

	t.Run("invalid_string", func(t *testing.T) {
		res, _ := Valid([]byte("9999"))

		assert.Equal(t, false, res)
	})
}
