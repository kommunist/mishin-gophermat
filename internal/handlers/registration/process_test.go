package registration

import (
	"bytes"
	"context"
	"encoding/json"
	"mishin-gophermat/internal/config"
	"net/http"
	"net/http/httptest"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestProcess(t *testing.T) {
	t.Run("Start_POST_to_create_user_in_db", func(t *testing.T) {
		// создали конфиг и стор
		c := config.MakeConfig()
		stor := NewMockAbstrStorage(gomock.NewController(t))

		// заинитили хендлер
		h := InitHandler(c, stor)

		// подготовили данные для запроса и сам запрос
		inputJSON, _ := json.Marshal(requestItem{Login: "Login", Password: "Password"})
		ctx := context.Background()
		request :=
			httptest.NewRequest(
				http.MethodPost,
				"/api/user/register",
				bytes.NewReader(inputJSON),
			).WithContext(ctx)

		// ожидаем, что в базу будет такой поход
		stor.EXPECT().CreateUser(ctx, "Login", "Password")

		// Делаем запрос
		w := httptest.NewRecorder()
		h.Process(w, request)
		res := w.Result()
		defer res.Body.Close()

		// Проверяем статус ответа
		assert.Equal(t, http.StatusOK, res.StatusCode)
		// Проверяем куку
		cookie := res.Cookies()[0]
		assert.Equal(t, "jwt", cookie.Name)
		assert.Greater(t, len(cookie.Value), 0)

		// Проверяем, что хедер совпадает с кукой
		header := res.Header.Get("Authorization")
		assert.Equal(t, "BEARER "+cookie.Value, header)

	})
}
