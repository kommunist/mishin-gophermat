package postorders

import (
	"bytes"
	"context"
	"mishin-gophermat/internal/secure"
	"net/http"
	"net/http/httptest"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// func GetLoginLenin(ctx context.Context) (jwt.Token, map[string]any, error) {
// 	return nil, map[string]any{"login": "lenin"}, nil
// }

func TestProcess(t *testing.T) {

	t.Run("create_order_202", func(t *testing.T) {
		acrChan := make(chan string, 5)

		// создали стор
		stor := NewMockOrderCreator(gomock.NewController(t))

		// заинитили хендлер
		h := InitHandler(stor, acrChan)

		//готовим запрос
		ctx := context.WithValue(context.Background(), secure.UserLoginKey, "lenin")
		request :=
			httptest.NewRequest(
				http.MethodPost,
				"/api/user/orders",
				bytes.NewReader([]byte("98265820")),
			).WithContext(ctx)

		// ожидаем, что в базу будет такой поход для поиска
		stor.EXPECT().OrderByNumberGet(ctx, "98265820")
		// ожидаем, что в базе будет создан заказ
		stor.EXPECT().OrderCreate(ctx, "98265820", "lenin")

		// Делаем запрос
		w := httptest.NewRecorder()
		h.Process(w, request)
		res := w.Result()
		defer res.Body.Close()

		// Проверяем статус ответа
		assert.Equal(t, http.StatusAccepted, res.StatusCode, "response status must be 202") // 202
	})

	t.Run("already_exist_order_with_same_author_200", func(t *testing.T) {
		acrChan := make(chan string, 5)
		// создали стор
		stor := NewMockOrderCreator(gomock.NewController(t))

		// заинитили хендлер
		h := InitHandler(stor, acrChan)

		//готовим запрос
		ctx := context.WithValue(context.Background(), secure.UserLoginKey, "lenin")
		request :=
			httptest.NewRequest(
				http.MethodPost,
				"/api/user/orders",
				bytes.NewReader([]byte("98265820")),
			).WithContext(ctx)

		// ожидаем, что в базу будет такой поход для поиска
		stor.EXPECT().OrderByNumberGet(ctx, "98265820").Return("lenin", nil)

		// Делаем запрос
		w := httptest.NewRecorder()
		h.Process(w, request)
		res := w.Result()
		defer res.Body.Close()

		// Проверяем статус ответа
		assert.Equal(t, http.StatusOK, res.StatusCode, "response status must be 200") // 200
	})

	t.Run("already_exist_order_with_another_author_409", func(t *testing.T) {
		acrChan := make(chan string, 5)
		// создали стор
		stor := NewMockOrderCreator(gomock.NewController(t))

		// заинитили хендлер
		h := InitHandler(stor, acrChan)

		//готовим запрос
		ctx := context.WithValue(context.Background(), secure.UserLoginKey, "lenin")
		request :=
			httptest.NewRequest(
				http.MethodPost,
				"/api/user/orders",
				bytes.NewReader([]byte("98265820")),
			).WithContext(ctx)

		// ожидаем, что в базу будет такой поход для поиска
		stor.EXPECT().OrderByNumberGet(ctx, "98265820").Return("bronstein", nil)

		// Делаем запрос
		w := httptest.NewRecorder()
		h.Process(w, request)
		res := w.Result()
		defer res.Body.Close()

		// Проверяем статус ответа
		assert.Equal(t, http.StatusConflict, res.StatusCode, "response status must be 409") // 409
	})

	t.Run("invalid_data_in_input_422", func(t *testing.T) {
		acrChan := make(chan string, 5)

		// создали стор
		stor := NewMockOrderCreator(gomock.NewController(t))

		// заинитили хендлер
		h := InitHandler(stor, acrChan)

		//готовим запрос
		ctx := context.WithValue(context.Background(), secure.UserLoginKey, "lenin")
		request :=
			httptest.NewRequest(
				http.MethodPost,
				"/api/user/orders",
				bytes.NewReader([]byte("1111")),
			).WithContext(ctx)

		// ожидаем, что в базу будет такой поход для поиска
		stor.EXPECT().OrderByNumberGet(ctx, "1111").Times(0)
		// ожидаем, что в базе будет создан заказ
		stor.EXPECT().OrderCreate(ctx, "1111", "lenin").Times(0)

		// Делаем запрос
		w := httptest.NewRecorder()
		h.Process(w, request)
		res := w.Result()
		defer res.Body.Close()

		// Проверяем статус ответа
		assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode, "response status must be 422") // 422
	})

	// t.Run("unauthorize_401", func(t *testing.T) {
	// 	acrChan := make(chan string, 5)

	// 	// создали стор
	// 	stor := NewMockOrderCreator(gomock.NewController(t))

	// 	// заинитили хендлер
	// 	h := InitHandler(stor, acrChan)
	// 	// h.GetLogin = GetLoginLenin - специально выключено, чтобы было видно, что не авторизовываем

	// 	//готовим запрос
	// 	ctx := context.Background()
	// 	request :=
	// 		httptest.NewRequest(
	// 			http.MethodPost,
	// 			"/api/user/orders",
	// 			bytes.NewReader([]byte("98265820")),
	// 		).WithContext(ctx)

	// 	// ожидаем, что в базу будет такой поход для поиска
	// 	stor.EXPECT().OrderByNumberGet(ctx, "98265820").Times(0)
	// 	// ожидаем, что в базе будет создан заказ
	// 	stor.EXPECT().OrderCreate(ctx, "98265820", "lenin").Times(0)

	// 	// Делаем запрос
	// 	w := httptest.NewRecorder()
	// 	h.Process(w, request)
	// 	res := w.Result()
	// 	defer res.Body.Close()

	// 	// Проверяем статус ответа
	// 	assert.Equal(t, http.StatusUnauthorized, res.StatusCode, "response status must be 401") // 401
	// })

}
