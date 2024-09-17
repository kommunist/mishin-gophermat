package accrual

import (
	"context"
	"log/slog"
	"time"
)

func (acr *Accrual) InitWorkers(inp chan string) {
	for i := 1; i < 5; i++ { // пусть будет 5 рутин для начала
		go acr.slave(inp)
	}
}

func (acr *Accrual) slave(inp chan string) {
	for num := range inp {
		founded, err := acr.process(num)
		if err != nil {
			slog.Error("error when process number from channel", "err", err)
			continue
		}
		if !founded {
			time.Sleep(5 * time.Second) // маленько подождем, чтобы не перегружать
			inp <- num                  // положим обратно
		}
	}
}

// - `REGISTERED` — заказ зарегистрирован, но не начисление не рассчитано;
// - `INVALID` — заказ не принят к расчёту, и вознаграждение не будет начислено;
// - `PROCESSING` — расчёт начисления в процессе;
// - `PROCESSED` — расчёт начисления окончен;
func (acr *Accrual) process(number string) (bool, error) { // repeat?(true/false)
	status, accrual, err := getOrderData(number)
	if err != nil {
		slog.Error("Error when get data from accrual", "err", err)
		return false, nil // будем считать, что достаточно того, что написали в логи
	}

	if status != "INVALID" && status != "PROCESSED" {
		// попробовать снова
		return true, nil
	}

	slog.Info("Try to update order", "number", number, "value", accrual, "status", status)

	err = acr.DB.UpdateOrderStatusAndValue(context.Background(), number, status, accrual)
	if err != nil {
		slog.Error("Error when update order in db", "err", err)
		return false, nil // будем считать, что достаточно того, что написали в логи. Возможно надо уходить на нвоый круг
	}

	return false, nil
}
