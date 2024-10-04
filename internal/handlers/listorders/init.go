package listorders

import (
	"context"
	"mishin-gophermat/internal/models"
)

type OrdersGetter interface {
	OrdersGet(ctx context.Context, login string) (data []models.Order, err error)
}

type ListOrdersHandler struct {
	DB OrdersGetter
}

func InitHandler(db OrdersGetter) ListOrdersHandler {
	return ListOrdersHandler{
		DB: db,
	}
}
