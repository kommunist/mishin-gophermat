package listwithdrawns

import (
	"context"
	"mishin-gophermat/internal/models"
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
		stor := NewMockWithdrawnsGetter(gomock.NewController(t))

		// заинитили хендлер
		h := InitHandler(stor)
		h.GetLogin = GetLoginLenin

		//готовим запрос
		ctx := context.Background()
		request :=
			httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil).WithContext(ctx)

		// ожидаем, что в базу будет такой поход для поиска
		stor.EXPECT().WithdrawnsGet(ctx, "lenin").Return(
			[]models.Withdrawn{
				{
					Number:      "number",
					Value:       500.0,
					ProcessedAt: "2021",
				},
			}, nil,
		)

		// Делаем запрос
		w := httptest.NewRecorder()
		h.Process(w, request)
		res := w.Result()
		defer res.Body.Close()

		// Проверяем статус ответа
		assert.Equal(t, http.StatusOK, res.StatusCode, "response status must be 200") // 200
	})

	t.Run("when_no_data_204", func(t *testing.T) {

		// создали стор
		stor := NewMockWithdrawnsGetter(gomock.NewController(t))

		// заинитили хендлер
		h := InitHandler(stor)
		h.GetLogin = GetLoginLenin

		//готовим запрос
		ctx := context.Background()
		request :=
			httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil).WithContext(ctx)

		// ожидаем, что в базу будет такой поход для поиска
		stor.EXPECT().WithdrawnsGet(ctx, "lenin").Return(nil, nil)

		// Делаем запрос
		w := httptest.NewRecorder()
		h.Process(w, request)
		res := w.Result()
		defer res.Body.Close()

		// Проверяем статус ответа
		assert.Equal(t, http.StatusNoContent, res.StatusCode, "response status must be 204") // 204
	})

	t.Run("when_unanauthorize_401", func(t *testing.T) {

		// создали стор
		stor := NewMockWithdrawnsGetter(gomock.NewController(t))

		// заинитили хендлер
		h := InitHandler(stor)
		// h.GetLogin = GetLoginLenin - специально выключено, чтобы было видно, что не авторизовываем

		//готовим запрос
		ctx := context.Background()
		request := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil).WithContext(ctx)

		// ожидаем, что в базу будет такой поход для поиска
		stor.EXPECT().WithdrawnsGet(ctx, "lenin").Times(0)

		// Делаем запрос
		w := httptest.NewRecorder()
		h.Process(w, request)
		res := w.Result()
		defer res.Body.Close()

		// Проверяем статус ответа
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode, "response status must be 401") // 401
	})
}
