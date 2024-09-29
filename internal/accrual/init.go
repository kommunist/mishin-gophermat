package accrual

import (
	"context"
	"sync"
)

type OrderUpdater interface {
	OrderUpdate(ctx context.Context, number string, status string, value float64) error
}

type Accrual struct {
	DB      OrderUpdater
	URI     string
	wg      sync.WaitGroup
	AcrChan chan string // канал для общения с воркерами
}

func InitAccrual(db OrderUpdater, URI string) Accrual {
	return Accrual{
		DB:      db,
		URI:     URI,
		AcrChan: make(chan string, 5),
	}
}

func (acr *Accrual) InitWorkers() {
	acr.wg.Add(5)
	for i := 1; i < 5; i++ { // пусть будет 5 рутин для начала
		go acr.slave(acr.AcrChan)
	}
}

func (acr *Accrual) FinishWorkers() {
	close(acr.AcrChan)

	acr.wg.Wait()
}
