package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"inventory-query-module/internal/bootstrap/config"
	"inventory-query-module/internal/bootstrap/module"

	configx "github.com/iamKienb/go-core/config"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type App struct {
	logger *slog.Logger
	server *http.Server
	infra  *module.InfraModule
}

func NewApp(logger *slog.Logger) *App {
	return &App{logger: logger}
}

func (a *App) Start(ctx context.Context) error {
	cfg, err := configx.Loader[config.InventoryQueryConfig]()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	if cfg == nil || cfg.Server.GrpcPort == 0 {
		return fmt.Errorf("config is empty: check your .env file path")
	}

	infra, err := module.NewInfraModule(cfg)
	if err != nil {
		return fmt.Errorf("infra: %w", err)
	}
	a.infra = infra

	application := module.NewApplicationModule(infra)
	adapter := module.NewAdapterModule(application, a.logger)

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

	a.logger.Info("starting inventory query", slog.Int("port", cfg.Server.GrpcPort))
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

	if a.infra != nil && a.infra.ESService != nil {
		if err := a.infra.ESService.Close(ctx); err != nil {
			return fmt.Errorf("close elasticsearch: %w", err)
		}
	}

	a.logger.Info("app stopped cleanly")
	return nil
}
