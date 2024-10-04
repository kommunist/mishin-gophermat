package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// подсмотерел shutdown в доках и чуть изменил под себя
func WaitFinish(app *App) {
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint // ждем сигнал о выключении

		slog.Info("Start shutdown!!!!!")
		app.Finish(false)
	}()

}

func (app *App) Finish(failed bool) {
	app.Acr.FinishWorkers()

	if err := app.srv.Shutdown(context.Background()); err != nil {
		slog.Warn("Error when sutdown server", "err", err)
	}

	if failed {
		os.Exit(1)
	} else {
		close(app.FinishChan)
	}
}
