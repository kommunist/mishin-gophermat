package postorders

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stretchr/testify/assert"
)

func GetLoginLenin(ctx context.Context) (jwt.Token, map[string]interface{}, error) {
	return nil, map[string]interface{}{"login": "lenin"}, nil
}

func TestProcess(t *testing.T) {

	t.Run("create_order_202", func(t *testing.T) {

		// создали стор
		stor := NewMockAbstrStorage(gomock.NewController(t))

		// заинитили хендлер
		h := InitHandler(stor)
		h.GetLogin = GetLoginLenin

		//готовим запрос
		ctx := context.Background()
		request :=
			httptest.NewRequest(
				http.MethodPost,
				"/api/user/orders",
				bytes.NewReader([]byte("98265820")),
			).WithContext(ctx)

		// ожидаем, что в базу будет такой поход для поиска
		stor.EXPECT().SelectOrderByNumber(ctx, "98265820")
		// ожидаем, что в базе будет создан заказ
		stor.EXPECT().CreateOrder(ctx, "98265820", "lenin")

		// Делаем запрос
		w := httptest.NewRecorder()
		h.Process(w, request)
		res := w.Result()
		defer res.Body.Close()

		// Проверяем статус ответа
		assert.Equal(t, http.StatusAccepted, res.StatusCode) // 202
	})

	t.Run("already_exist_order_with_same_author_200", func(t *testing.T) {
		// создали стор
		stor := NewMockAbstrStorage(gomock.NewController(t))

		// заинитили хендлер
		h := InitHandler(stor)
		h.GetLogin = GetLoginLenin

		//готовим запрос
		ctx := context.Background()
		request :=
			httptest.NewRequest(
				http.MethodPost,
				"/api/user/orders",
				bytes.NewReader([]byte("98265820")),
			).WithContext(ctx)

		// ожидаем, что в базу будет такой поход для поиска
		stor.EXPECT().SelectOrderByNumber(ctx, "98265820").Return(
			map[string]interface{}{"userLogin": "lenin"}, nil,
		)

		// Делаем запрос
		w := httptest.NewRecorder()
		h.Process(w, request)
		res := w.Result()
		defer res.Body.Close()

		// Проверяем статус ответа
		assert.Equal(t, http.StatusOK, res.StatusCode) // 200
	})

	t.Run("already_exist_order_with_another_author_409", func(t *testing.T) {
		// создали стор
		stor := NewMockAbstrStorage(gomock.NewController(t))

		// заинитили хендлер
		h := InitHandler(stor)
		h.GetLogin = GetLoginLenin

		//готовим запрос
		ctx := context.Background()
		request :=
			httptest.NewRequest(
				http.MethodPost,
				"/api/user/orders",
				bytes.NewReader([]byte("98265820")),
			).WithContext(ctx)

		// ожидаем, что в базу будет такой поход для поиска
		stor.EXPECT().SelectOrderByNumber(ctx, "98265820").Return(
			map[string]interface{}{"userLogin": "bronstein"}, nil,
		)

		// Делаем запрос
		w := httptest.NewRecorder()
		h.Process(w, request)
		res := w.Result()
		defer res.Body.Close()

		// Проверяем статус ответа
		assert.Equal(t, http.StatusConflict, res.StatusCode) // 409
	})

	t.Run("invalid data in input", func(t *testing.T) {

		// создали стор
		stor := NewMockAbstrStorage(gomock.NewController(t))

		// заинитили хендлер
		h := InitHandler(stor)
		h.GetLogin = GetLoginLenin

		//готовим запрос
		ctx := context.Background()
		request :=
			httptest.NewRequest(
				http.MethodPost,
				"/api/user/orders",
				bytes.NewReader([]byte("1111")),
			).WithContext(ctx)

		// ожидаем, что в базу будет такой поход для поиска
		stor.EXPECT().SelectOrderByNumber(ctx, "1111").Times(0)
		// ожидаем, что в базе будет создан заказ
		stor.EXPECT().CreateOrder(ctx, "1111", "lenin").Times(0)

		// Делаем запрос
		w := httptest.NewRecorder()
		h.Process(w, request)
		res := w.Result()
		defer res.Body.Close()

		// Проверяем статус ответа
		assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode) // 422
	})

	t.Run("anauthorize", func(t *testing.T) {

		// создали стор
		stor := NewMockAbstrStorage(gomock.NewController(t))

		// заинитили хендлер
		h := InitHandler(stor)
		// h.GetLogin = GetLoginLenin - специально выключено, чтобы было видно, что не авторизовываем

		//готовим запрос
		ctx := context.Background()
		request :=
			httptest.NewRequest(
				http.MethodPost,
				"/api/user/orders",
				bytes.NewReader([]byte("98265820")),
			).WithContext(ctx)

		// ожидаем, что в базу будет такой поход для поиска
		stor.EXPECT().SelectOrderByNumber(ctx, "98265820").Times(0)
		// ожидаем, что в базе будет создан заказ
		stor.EXPECT().CreateOrder(ctx, "98265820", "lenin").Times(0)

		// Делаем запрос
		w := httptest.NewRecorder()
		h.Process(w, request)
		res := w.Result()
		defer res.Body.Close()

		// Проверяем статус ответа
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode) // 422
	})

}
