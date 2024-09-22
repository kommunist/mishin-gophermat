package balance

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stretchr/testify/assert"
)

func GetLoginLenin(ctx context.Context) (jwt.Token, map[string]any, error) {
	return nil, map[string]any{"login": "lenin"}, nil
}

func TestProcess(t *testing.T) {

	t.Run("correct_return_list_200", func(t *testing.T) {

		// создали стор
		stor := NewMockBalanceGetter(gomock.NewController(t))

		// заинитили хендлер
		h := InitHandler(stor)
		h.GetLogin = GetLoginLenin

		//готовим запрос
		ctx := context.Background()
		request :=
			httptest.NewRequest(http.MethodGet, "/api/user/balance", nil).WithContext(ctx)

		// ожидаем, что в базу будет такой поход для поиска
		stor.EXPECT().BalanceGet(ctx, "lenin").Return(
			500.0, 60.0, nil,
		)

		// Делаем запрос
		w := httptest.NewRecorder()
		h.Process(w, request)
		res := w.Result()
		defer res.Body.Close()

		body, _ := io.ReadAll(res.Body)
		resp := response{}
		json.Unmarshal(body, &resp)

		assert.Equal(t, 500.0, resp.Current, "current balance must be eq to result from db")
		assert.Equal(t, 60.0, resp.Withdrawn, "current size of withdrawns must be eq to result from db")

		// Проверяем статус ответа
		assert.Equal(t, http.StatusOK, res.StatusCode, "response status must be 200") // 200
	})

	t.Run("when_without_authorize_401", func(t *testing.T) {

		// создали стор
		stor := NewMockBalanceGetter(gomock.NewController(t))

		// заинитили хендлер
		h := InitHandler(stor)
		// h.GetLogin = GetLoginLenin

		//готовим запрос
		ctx := context.Background()
		request :=
			httptest.NewRequest(http.MethodGet, "/api/user/balance", nil).WithContext(ctx)

		// ожидаем, что в базу будет такой поход для поиска
		stor.EXPECT().BalanceGet(ctx, "lenin").Times(0)

		// Делаем запрос
		w := httptest.NewRecorder()
		h.Process(w, request)
		res := w.Result()
		defer res.Body.Close()

		// Проверяем статус ответа
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode, "response status must be 401") // 401
	})
}
