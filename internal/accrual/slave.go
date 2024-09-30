package accrual

import (
	"context"
	"log/slog"
	"time"
)

func (acr *Accrual) slave(inp chan string) {
	defer acr.wg.Done()

	for num := range inp {
		repeat := acr.process(num)
		if repeat {
			time.Sleep(5 * time.Second) // маленько подождем, чтобы не перегружать
			inp <- num                  // положим обратно
		}
	}
}

// - `REGISTERED` — заказ зарегистрирован, но не начисление не рассчитано;
// - `INVALID` — заказ не принят к расчёту, и вознаграждение не будет начислено;
// - `PROCESSING` — расчёт начисления в процессе;
// - `PROCESSED` — расчёт начисления окончен;
func (acr *Accrual) process(number string) bool { // repeat?(true/false)
	status, accrual, err := acr.getOrderData(number)
	if err != nil {
		slog.Error("Error when get data from accrual", "err", err)
		return false // будем считать, что достаточно того, что написали в логи
	}

	if status != "INVALID" && status != "PROCESSED" {
		// попробовать снова
		return true
	}

	slog.Info("Try to update order", "number", number, "value", accrual, "status", status)

	err = acr.DB.OrderUpdate(context.Background(), number, status, accrual)
	if err != nil {
		slog.Error("Error when update order in db", "err", err)
		return false // будем считать, что достаточно того, что написали в логи. Возможно надо уходить на нвоый круг
	}

	return false
}
