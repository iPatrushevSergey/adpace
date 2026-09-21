package bootstrap

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/port"
)

type App struct {
	Server          *http.Server
	Log             port.Logger
	ShutdownTimeout time.Duration
	TLSCertFile     string
	TLSKeyFile      string
}

func (a *App) Start() <-chan error {
	errChan := make(chan error, 1)

	go func() {
		a.Log.Info(context.Background(), "server listening", "address", a.Server.Addr, "tls", a.TLSCertFile != "")
		var err error
		if a.TLSCertFile != "" && a.TLSKeyFile != "" {
			err = a.Server.ListenAndServeTLS(a.TLSCertFile, a.TLSKeyFile)
		} else {
			err = a.Server.ListenAndServe()
		}
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	return errChan
}

func (a *App) Shutdown(ctx context.Context) error {
	a.Log.Info(context.Background(), "stopping server...")
	if err := a.Server.Shutdown(ctx); err != nil {
		return err
	}
	a.Log.Info(context.Background(), "server stopped gracefully")
	return nil
}
