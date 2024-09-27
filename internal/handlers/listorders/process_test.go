package listorders

import (
	"context"
	"mishin-gophermat/internal/models"
	"mishin-gophermat/internal/secure"
	"net/http"
	"net/http/httptest"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestProcess(t *testing.T) {

	t.Run("correct_return_list_200", func(t *testing.T) {

		// создали стор
		stor := NewMockOrdersGetter(gomock.NewController(t))

		// заинитили хендлер
		h := InitHandler(stor)
		ctx := context.WithValue(context.Background(), secure.UserLoginKey, "lenin")

		//готовим запрос
		request :=
			httptest.NewRequest(http.MethodGet, "/api/user/orders", nil).WithContext(ctx)

		// ожидаем, что в базу будет такой поход для поиска
		stor.EXPECT().OrdersGet(ctx, "lenin").Return(
			[]models.Order{
				{
					Number:     "123",
					Status:     "NEW",
					Value:      500.0,
					UploadedAt: "2021",
				},
			}, nil,
		)

		// Делаем запрос
		w := httptest.NewRecorder()
		h.Process(w, request)
		res := w.Result()
		defer res.Body.Close()

		// Проверяем статус ответа
		assert.Equal(t, http.StatusOK, res.StatusCode, "response status must be 200")
	})

	t.Run("when_no_data_204", func(t *testing.T) {

		// создали стор
		stor := NewMockOrdersGetter(gomock.NewController(t))

		// заинитили хендлер
		h := InitHandler(stor)

		//готовим запрос
		ctx := context.WithValue(context.Background(), secure.UserLoginKey, "lenin")
		request :=
			httptest.NewRequest(http.MethodGet, "/api/user/orders", nil).WithContext(ctx)

		// ожидаем, что в базу будет такой поход для поиска
		stor.EXPECT().OrdersGet(ctx, "lenin").Return(nil, nil)

		// Делаем запрос
		w := httptest.NewRecorder()
		h.Process(w, request)
		res := w.Result()
		defer res.Body.Close()

		// Проверяем статус ответа
		assert.Equal(t, http.StatusNoContent, res.StatusCode, "response status must be 204")
	})

	t.Run("when_without_login_in_context_500", func(t *testing.T) {

		// создали стор
		stor := NewMockOrdersGetter(gomock.NewController(t))

		// заинитили хендлер
		h := InitHandler(stor)

		//готовим запрос
		ctx := context.Background()
		request := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil).WithContext(ctx)

		// ожидаем, что в базу будет такой поход для поиска
		stor.EXPECT().OrdersGet(ctx, "lenin").Times(0)

		// Делаем запрос
		w := httptest.NewRecorder()
		h.Process(w, request)
		res := w.Result()
		defer res.Body.Close()

		// Проверяем статус ответа
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode, "response status must be 500")
	})
}
