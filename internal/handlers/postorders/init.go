package postorders

import (
	"context"
)

type OrderCreator interface {
	OrderByNumberGet(ctx context.Context, numer string) (login string, err error)
	OrderCreate(ctx context.Context, number string, userLogin string) error
}

type PostOrdersHandler struct {
	DB      OrderCreator
	acrChan chan string
}

func InitHandler(db OrderCreator, acrChan chan string) PostOrdersHandler {
	return PostOrdersHandler{
		DB:      db,
		acrChan: acrChan,
	}
}
