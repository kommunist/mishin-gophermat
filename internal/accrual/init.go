package accrual

import "context"

type AbstrStorage interface {
	UpdateOrderStatusAndValue(ctx context.Context, number string, status string, value float64) error
}

type Accrual struct {
	DB  AbstrStorage
	URI string
}

func InitAccrual(db AbstrStorage, URI string) Accrual {
	return Accrual{DB: db, URI: URI}
}
