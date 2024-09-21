package accrual

import "context"

type OrderUpdater interface {
	OrderUpdate(ctx context.Context, number string, status string, value float64) error
}

type Accrual struct {
	DB  OrderUpdater
	URI string
}

func InitAccrual(db OrderUpdater, URI string) Accrual {
	return Accrual{DB: db, URI: URI}
}
