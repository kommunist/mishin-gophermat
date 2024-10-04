package login

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type FakeTrueChecker struct{} // для проверок положительного сценария

func (c *FakeTrueChecker) PassCheck(pass string, hashed string) bool {
	return true
}

type FakeFalseChecker struct{} // для проверок отрицательного сценария

func (c *FakeFalseChecker) PassCheck(pass string, hashed string) bool {
	return false
}

func TestProcess(t *testing.T) {
	t.Run("login_user_happy_path_200", func(t *testing.T) {
		// создали конфиг и стор
		stor := NewMockUserGetter(gomock.NewController(t))

		// заинитили хендлер
		h := InitHandler(stor)
		h.checker = &FakeTrueChecker{}

		// подготовили данные для запроса и сам запрос
		inputJSON, _ := json.Marshal(requestItem{Login: "Login", Password: "Password"})
		ctx := context.Background()
		request :=
			httptest.NewRequest(
				http.MethodPost,
				"/api/user/login",
				bytes.NewReader(inputJSON),
			).WithContext(ctx)

		// ожидаем, что в базу будет такой поход
		stor.EXPECT().UserGet(ctx, "Login").Return("hashed", nil)

		// Делаем запрос
		w := httptest.NewRecorder()
		h.Process(w, request)
		res := w.Result()
		defer res.Body.Close()

		// Проверяем статус ответа
		assert.Equal(t, http.StatusOK, res.StatusCode, "response status must be 200") // 200
		// Проверяем куку
		cookie := res.Cookies()[0]
		assert.Equal(t, "jwt", cookie.Name, "response has cookie with name jwt")
		assert.Greater(t, len(cookie.Value), 0, "length of this cookie is positive")

		// Проверяем, что хедер совпадает с кукой
		header := res.Header.Get("Authorization")
		assert.Equal(t, "BEARER "+cookie.Value, header, "authorization cookie must has same as cookie value")
	})

	t.Run("login_user_in_db_incorrect_format_400", func(t *testing.T) {
		// создали конфиг и стор
		stor := NewMockUserGetter(gomock.NewController(t))

		// заинитили хендлер
		h := InitHandler(stor)
		h.checker = &FakeTrueChecker{}

		// подготовили данные для запроса и сам запрос
		ctx := context.Background()
		request :=
			httptest.NewRequest(
				http.MethodPost,
				"/api/user/login",
				bytes.NewReader([]byte("vupsenpupsen")), // а запрос то кривой
			).WithContext(ctx)

		// т.е. ожидаем, что запроса в базу не будет
		stor.EXPECT().UserGet(ctx, "Login").Times(0)

		// Делаем запрос
		w := httptest.NewRecorder()
		h.Process(w, request)
		res := w.Result()
		defer res.Body.Close()

		// Проверяем статус ответа
		assert.Equal(t, http.StatusBadRequest, res.StatusCode, "response status must be 400") // 400
	})

	t.Run("login_user_with_incorrect_password", func(t *testing.T) {
		// создали конфиг и стор
		stor := NewMockUserGetter(gomock.NewController(t))

		// заинитили хендлер
		h := InitHandler(stor)
		h.checker = &FakeFalseChecker{}

		// подготовили данные для запроса и сам запрос
		inputJSON, _ := json.Marshal(requestItem{Login: "Login", Password: "Password"})
		ctx := context.Background()
		request :=
			httptest.NewRequest(
				http.MethodPost,
				"/api/user/login",
				bytes.NewReader(inputJSON),
			).WithContext(ctx)

		// ожидаем, что в базу будет такой поход
		stor.EXPECT().UserGet(ctx, "Login").Return("hashed", nil)

		// Делаем запрос
		w := httptest.NewRecorder()
		h.Process(w, request)
		res := w.Result()
		defer res.Body.Close()

		// Проверяем статус ответа
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode, "response status must be 401") // 401
		// Проверяем отсутствие куки
		cookies := res.Cookies()
		assert.Equal(t, 0, len(cookies), "response must be without cookies")

		// Проверяем, что хедера нет
		header := res.Header.Get("Authorization")
		assert.Equal(t, "", header, "response must be without authorization header")

	})

}
