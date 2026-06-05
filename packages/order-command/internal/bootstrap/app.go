package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"order-command-module/internal/bootstrap/config"
	"order-command-module/internal/bootstrap/module"
	"order-command-module/internal/infra/temporal"
	"order-command-module/internal/infra/temporal/activity"

	configx "github.com/iamKienb/go-core/config"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type App struct {
	logger *slog.Logger
	server *http.Server
	infra  *module.InfraModule
	worker *temporal.Worker
}

func NewApp(logger *slog.Logger) *App {
	return &App{logger: logger}
}

func (a *App) Start(ctx context.Context) error {
	cfg, err := configx.Loader[config.OrderCommandConfig]()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	if cfg == nil || cfg.Server.GrpcPort == 0 {
		return fmt.Errorf("config is empty: check your .env file path")
	}

	infra, err := module.NewInfraModule(ctx, cfg)
	if err != nil {
		return fmt.Errorf("infra: %w", err)
	}
	a.infra = infra

	application := module.NewApplicationModule(infra)
	adapter := module.NewAdapterModule(application, a.logger)

	orderActivity := activity.NewOrderActivity(
		application.OrderService,
		infra.UserClient,
		infra.ProductClient,
		infra.InventoryClient,
	)
	a.worker = temporal.NewWorker(infra.TemporalClient, cfg.Temporal.OrderTaskQueue, temporal.Registry{
		OrderActivity: orderActivity,
	})
	if err := a.worker.Start(); err != nil {
		return fmt.Errorf("start temporal worker: %w", err)
	}

	a.server = &http.Server{
		Addr: ":" + strconv.Itoa(cfg.Server.GrpcPort),
		Handler: h2c.NewHandler(
			adapter.Mux,
			&http2.Server{},
		),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	a.logger.Info("starting", slog.Int("port", cfg.Server.GrpcPort))
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server: %w", err)
	}

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	a.logger.Info("shutting down")

	if a.server != nil {
		if err := a.server.Shutdown(ctx); err != nil {
			return fmt.Errorf("shutdown server: %w", err)
		}
	}

	if a.worker != nil {
		a.worker.Stop()
	}

	if a.infra != nil && a.infra.PGService != nil {
		a.logger.Info("closing postgres connection")
		a.infra.PGService.Close()
	}

	a.logger.Info("app stopped cleanly")
	return nil
}
