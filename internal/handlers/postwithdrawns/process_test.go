package postwithdrawns

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/likexian/gokit/assert"
)

func GetLoginLenin(ctx context.Context) (jwt.Token, map[string]interface{}, error) {
	return nil, map[string]interface{}{"login": "lenin"}, nil
}

func TestProcess(t *testing.T) {
	t.Run("create_withdrawn_when_balance_is_high_200", func(t *testing.T) {

		// создали стор
		stor := NewMockAbstrStorage(gomock.NewController(t))

		// инитим хендлер
		h := InitHandler(stor)
		h.GetLogin = GetLoginLenin

		data, _ := json.Marshal(request{Number: "new_number", Value: 123})

		// готовим запрос
		ctx := context.Background()
		request :=
			httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewReader(data)).WithContext(ctx)

		// ожидаем, что в базу будет такой поход для баланса
		stor.EXPECT().SelectBalanceByLogin(ctx, "lenin").Return(500.0, 0.0, nil)
		// ожидаем, что в базе будет создан заказ
		stor.EXPECT().CreateWithdrawn(ctx, "lenin", "new_number", 123.0).Return(nil)

		w := httptest.NewRecorder()
		h.Process(w, request)
		res := w.Result()
		defer res.Body.Close()

		// Проверяем статус ответа
		assert.Equal(t, http.StatusOK, res.StatusCode) // 200
	})

	t.Run("create_withdrawn_when_balance_is_low_402", func(t *testing.T) {

		// создали стор
		stor := NewMockAbstrStorage(gomock.NewController(t))

		// инитим хендлер
		h := InitHandler(stor)
		h.GetLogin = GetLoginLenin

		data, _ := json.Marshal(request{Number: "new_number", Value: 123})

		// готовим запрос
		ctx := context.Background()
		request :=
			httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewReader(data)).WithContext(ctx)

		// ожидаем, что в базу будет такой поход для баланса и мало денег в балансе
		stor.EXPECT().SelectBalanceByLogin(ctx, "lenin").Return(100.0, 0.0, nil)
		// ожидаем, что в базе будет создан заказ
		stor.EXPECT().CreateWithdrawn(ctx, "lenin", "new_number", 123.0).Times(0)

		w := httptest.NewRecorder()
		h.Process(w, request)
		res := w.Result()
		defer res.Body.Close()

		// Проверяем статус ответа
		assert.Equal(t, http.StatusPaymentRequired, res.StatusCode) // 402
	})

	t.Run("when_user_not_authorized", func(t *testing.T) {

		// создали стор
		stor := NewMockAbstrStorage(gomock.NewController(t))

		// инитим хендлер
		h := InitHandler(stor)
		// h.GetLogin = GetLoginLenin // специально выключено

		data, _ := json.Marshal(request{Number: "new_number", Value: 123})

		// готовим запрос
		ctx := context.Background()
		request :=
			httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewReader(data)).WithContext(ctx)

		// ожидаем, что в базу будет такой поход для баланса
		stor.EXPECT().SelectBalanceByLogin(ctx, "lenin").Times(0)
		// ожидаем, что в базе будет создан заказ
		stor.EXPECT().CreateWithdrawn(ctx, "lenin", "new_number", 123).Times(0)

		w := httptest.NewRecorder()
		h.Process(w, request)
		res := w.Result()
		defer res.Body.Close()

		// Проверяем статус ответа
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode) // 401
	})
}
