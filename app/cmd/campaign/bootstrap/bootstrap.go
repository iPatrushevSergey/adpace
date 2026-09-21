package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/iPatrushevSergey/adpace/app/internal/campaign/adapter/clock"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/adapter/generator"
	campaignpostgres "github.com/iPatrushevSergey/adpace/app/internal/campaign/adapter/repository/postgres"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/usecase"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/config"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/presentation/http/router"
	"github.com/iPatrushevSergey/adpace/app/internal/pkg/adapter/logger"
	"github.com/iPatrushevSergey/adpace/app/internal/pkg/adapter/repository/postgres"
	"github.com/spf13/pflag"
)

func Run() (*App, []func(), error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		if errors.Is(err, pflag.ErrHelp) {
			return nil, nil, nil
		}
		return nil, nil, fmt.Errorf("load config: %w", err)
	}

	log, err := logger.NewSlogLogger(cfg.Logger)
	if err != nil {
		return nil, nil, fmt.Errorf("init logger: %w", err)
	}

	log.Info(
		context.Background(),
		"starting campaign-manager",
		"address", cfg.Server.Address,
		"tls_configured", cfg.Server.TLSEnabled(),
	)

	pool, err := postgres.NewPool(context.Background(), cfg.DBPool)
	if err != nil {
		return nil, nil, fmt.Errorf("database pool: %w", err)
	}
	cleanups := []func(){func() { pool.Close() }}
	log.Info(context.Background(), "database connected")

	transactor := postgres.NewTransactor(pool)
	retryer := postgres.NewRetryer(
		postgres.WithMaxRetries(cfg.DBRetry.MaxRetries),
		postgres.WithExponentialBackoff(cfg.DBRetry.BaseDelay, cfg.DBRetry.MaxDelay),
	)

	advertiserRepo := campaignpostgres.NewAdvertiserRepo(transactor, retryer)
	idGen := generator.NewIDGenerator()
	clk := clock.NewRealClock()

	uc := usecase.AdvertiserUseCases{
		Create: usecase.NewCreateAdvertiser(advertiserRepo, idGen, clk),
		Get:    usecase.NewGetByIDAdvertiser(advertiserRepo),
		Patch:  usecase.NewPatchAdvertiser(advertiserRepo, clk),
		Put:    usecase.NewPutAdvertiser(advertiserRepo, clk),
		Delete: usecase.NewDeleteAdvertiser(advertiserRepo),
	}

	r := router.New(uc, log)

	srv := &http.Server{
		Addr:    cfg.Server.Address,
		Handler: r,
	}

	app := &App{
		Server:          srv,
		Log:             log,
		ShutdownTimeout: cfg.Server.ShutdownTimeout,
		TLSCertFile:     cfg.Server.CertFile,
		TLSKeyFile:      cfg.Server.KeyFile,
	}

	return app, cleanups, nil
}
