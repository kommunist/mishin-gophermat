package accrual

import "context"

type AbstrStorage interface {
	UpdateOrderStatusAndValue(ctx context.Context, number string, status string, value float64) error
}

type Accrual struct {
	DB AbstrStorage
}

func InitAccrual(db AbstrStorage) Accrual {
	return Accrual{DB: db}
}
