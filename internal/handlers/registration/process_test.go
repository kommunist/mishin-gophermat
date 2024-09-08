package registration

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mishin-gophermat/internal/errors/exist"
	"net/http"
	"net/http/httptest"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestProcess(t *testing.T) {
	t.Run("create_user_in_db_happy_path_200", func(t *testing.T) {
		// создали конфиг и стор
		stor := NewMockAbstrStorage(gomock.NewController(t))

		// заинитили хендлер
		h := InitHandler(stor)

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
		assert.Equal(t, http.StatusOK, res.StatusCode) // 200
		// Проверяем куку
		cookie := res.Cookies()[0]
		assert.Equal(t, "jwt", cookie.Name)
		assert.Greater(t, len(cookie.Value), 0)

		// Проверяем, что хедер совпадает с кукой
		header := res.Header.Get("Authorization")
		assert.Equal(t, "BEARER "+cookie.Value, header)

	})

	t.Run("create_user_in_db_incorrect_format_400", func(t *testing.T) {
		// создали конфиг и стор
		stor := NewMockAbstrStorage(gomock.NewController(t))

		// заинитили хендлер
		h := InitHandler(stor)

		// подготовили данные для запроса и сам запрос
		ctx := context.Background()
		request :=
			httptest.NewRequest(
				http.MethodPost,
				"/api/user/register",
				bytes.NewReader([]byte("vupsenpupsen")), // а запрос то кривой
			).WithContext(ctx)

		// т.е. ожидаем, что запроса в базу не будет
		stor.EXPECT().CreateUser(ctx, "Login", "Password").Times(0)

		// Делаем запрос
		w := httptest.NewRecorder()
		h.Process(w, request)
		res := w.Result()
		defer res.Body.Close()

		// Проверяем статус ответа
		assert.Equal(t, http.StatusBadRequest, res.StatusCode) // 400
	})

	t.Run("create_user_in_db_already_exist_409", func(t *testing.T) {
		// создали конфиг и стор
		stor := NewMockAbstrStorage(gomock.NewController(t))

		// заинитили хендлер
		h := InitHandler(stor)

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
		stor.EXPECT().CreateUser(
			ctx, "Login", "Password",
		).Return(exist.NewExistError(errors.New("qq")))

		// Делаем запрос
		w := httptest.NewRecorder()
		h.Process(w, request)
		res := w.Result()
		defer res.Body.Close()

		// Проверяем статус ответа
		assert.Equal(t, http.StatusConflict, res.StatusCode) // 409
		// Проверяем отсутствие куки
		cookies := res.Cookies()
		assert.Equal(t, 0, len(cookies))

		// Проверяем, что хедера нет
		header := res.Header.Get("Authorization")
		assert.Equal(t, "", header)

	})

}
