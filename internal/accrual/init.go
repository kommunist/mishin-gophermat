package accrual

import (
	"context"
	"sync"
)

type OrderUpdater interface {
	OrderUpdate(ctx context.Context, number string, status string, value float64) error
}

type Accrual struct {
	DB       OrderUpdater
	URI      string
	wg       sync.WaitGroup
	AcrChan  chan string // канал для общения с воркерами
	Count    int
	waitChan chan int // канал для ожидания
}

func InitAccrual(db OrderUpdater, URI string) Accrual {
	return Accrual{
		DB:       db,
		URI:      URI,
		AcrChan:  make(chan string, 5),
		Count:    5, // пока пусть будет 5 воркеров
		waitChan: make(chan int, 5),
	}
}

func (acr *Accrual) InitWorkers() {
	acr.wg.Add(5)
	for i := 0; i < acr.Count; i++ {
		go acr.slave()
	}
}

func (acr *Accrual) FinishWorkers() {
	close(acr.AcrChan)

	acr.wg.Wait()
}
