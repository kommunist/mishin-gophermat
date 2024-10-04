package accrual

import (
	"context"
	"log/slog"
	"time"
)

func (acr *Accrual) slave() {
	defer acr.wg.Done()
	closed := false

	for !closed {
		select {
		case num, ok := <-acr.AcrChan:
			if ok {
				acr.process(num)
			}
			closed = !ok
		case wait := <-acr.waitChan:
			time.Sleep(time.Duration(wait) * time.Second)
		}
	}
}

func (acr *Accrual) process(number string) {
	status, accrual, wait, err := acr.getOrderData(number)
	if err != nil {
		slog.Error("Error when get data from accrual", "err", err)
		acr.AcrChan <- number // будем прокручивать, чтобы не потерять
		return
	}

	if wait != 0 {
		for i := 0; i < acr.Count; i++ {
			acr.waitChan <- wait
		}
	}

	if status != "INVALID" && status != "PROCESSED" {
		acr.AcrChan <- number // попробовать снова
		return
	}

	err = acr.DB.OrderUpdate(context.Background(), number, status, accrual)
	if err != nil {
		slog.Error("Error when update order in db", "err", err)
		acr.AcrChan <- number // будем прокручивать, чтобы не потерять
		return
	}
}
