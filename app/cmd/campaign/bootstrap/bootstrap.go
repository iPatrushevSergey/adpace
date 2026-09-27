package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/adapter/clock"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/adapter/generator"
	campaignpostgres "github.com/iPatrushevSergey/adpace/app/internal/campaign/adapter/repository/postgres"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/port"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/usecase"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/config"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/presentation/http/router"
	"github.com/iPatrushevSergey/adpace/app/internal/pkg/adapter/logger"
	"github.com/iPatrushevSergey/adpace/app/internal/pkg/adapter/repository/postgres"
	"github.com/iPatrushevSergey/adpace/app/internal/pkg/adapter/retry"
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
	retryer := retry.NewRetryer(
		port.WithAttempts(cfg.DBRetry.Attempts),
		port.WithBackoffFunc(retry.ExponentialBackoff(cfg.DBRetry.BaseDelay, cfg.DBRetry.MaxDelay)),
	)

	getter := trmpgx.DefaultCtxGetter
	advertiserRepo := campaignpostgres.NewAdvertiserRepo(pool, getter)
	campaignRepo := campaignpostgres.NewCampaignRepo(pool, getter)
	idGen := generator.NewIDGenerator()
	clk := clock.NewRealClock()

	advUC := usecase.AdvertiserUseCases{
		Create: usecase.NewCreateAdvertiser(advertiserRepo, idGen, clk, retryer),
		Get:    usecase.NewGetByIDAdvertiser(advertiserRepo, retryer),
		Patch:  usecase.NewPatchAdvertiser(advertiserRepo, transactor, retryer, clk),
		Put:    usecase.NewPutAdvertiser(advertiserRepo, clk, retryer),
		Delete: usecase.NewDeleteAdvertiser(advertiserRepo, retryer),
	}

	campUC := usecase.CampaignUseCases{
		Create: usecase.NewCreateCampaign(campaignRepo, idGen, clk, retryer),
		Get:    usecase.NewGetByIDCampaign(campaignRepo, retryer),
		Patch:  usecase.NewPatchCampaign(campaignRepo, transactor, retryer, clk),
		Put:    usecase.NewPutCampaign(campaignRepo, clk, retryer),
		Delete: usecase.NewDeleteCampaign(campaignRepo, retryer),
		Pause:  usecase.NewPauseCampaign(campaignRepo, clk, retryer),
		Resume: usecase.NewResumeCampaign(campaignRepo, clk, retryer),
	}

	r := router.New(advUC, campUC, log)

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
