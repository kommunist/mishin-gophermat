package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/likexian/gokit/assert"
)

func testHandler(w http.ResponseWriter, r *http.Request) {}

func FakeLogin(ctx context.Context) (jwt.Token, map[string]any, error) {
	return nil, map[string]any{"login": "lenin"}, nil
}

func TestAuth(t *testing.T) {
	t.Run("correct_auth", func(t *testing.T) {
		a := InitAuth()
		a.GetLogin = FakeLogin

		nextHandler := http.HandlerFunc(testHandler)
		handlerToTest := a.Auth(nextHandler)

		request :=
			httptest.NewRequest(http.MethodGet, "/any", nil)

		w := httptest.NewRecorder()
		handlerToTest.ServeHTTP(w, request)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, res.StatusCode, http.StatusOK, "response status must be 200 with auth")
	})

	t.Run("incorrect_auth", func(t *testing.T) {
		a := InitAuth()
		// a.GetLogin = FakeLogin

		nextHandler := http.HandlerFunc(testHandler)
		handlerToTest := a.Auth(nextHandler)

		request :=
			httptest.NewRequest(http.MethodGet, "/any", nil)

		w := httptest.NewRecorder()
		handlerToTest.ServeHTTP(w, request)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, res.StatusCode, http.StatusUnauthorized, "response status must be 401 without auth")
	})
}
