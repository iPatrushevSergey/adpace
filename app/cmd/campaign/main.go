package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/iPatrushevSergey/adpace/app/cmd/campaign/bootstrap"
)

func main() {
	app, cleanups, err := bootstrap.Run()
	for _, cleanup := range cleanups {
		defer cleanup()
	}
	if err != nil {
		log.Fatalf("campaign: %v", err)
	}
	if app == nil {
		return
	}

	errChan := app.Start()

	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	var startErr error
	select {
	case <-signalCtx.Done():
	case startErr = <-errChan:
	}
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), app.ShutdownTimeout)
	defer cancel()

	shutdownErr := app.Shutdown(ctx)

	if startErr != nil {
		app.Log.Error(context.Background(), "server start error", "error", startErr)
	}
	if shutdownErr != nil {
		app.Log.Error(context.Background(), "shutdown error", "error", shutdownErr)
	}
	if startErr != nil || shutdownErr != nil {
		os.Exit(1)
	}
}
