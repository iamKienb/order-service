package module

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"order-command-module/internal/application/port"
	"order-command-module/internal/bootstrap/config"
	"order-command-module/internal/domain/order"
	"order-command-module/internal/infra/grpc_client"
	orderPg "order-command-module/internal/infra/postgres/order"
	outboxPg "order-command-module/internal/infra/postgres/outbox"
	"order-command-module/internal/infra/temporal/runner"

	pgx "github.com/iamKienb/go-core/postgres"
	"go.temporal.io/sdk/client"
)

type InfraModule struct {
	TemporalClient client.Client
	WorkflowRunner runner.Runner
	PGService      pgx.PGXService

	OrderRepo  order.Repository
	OutboxRepo port.OutboxRepository

	UserClient      port.UserClient
	ProductClient   port.ProductClient
	InventoryClient port.InventoryClient

	TxManager port.TxManager
}

func NewInfraModule(ctx context.Context, cfg *config.OrderCommandConfig) (*InfraModule, error) {
	pgService, err := pgx.New(cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}

	temporalClient, err := client.Dial(client.Options{
		HostPort:  cfg.Temporal.Address,
		Namespace: cfg.Temporal.Namespace,
	})
	if err != nil {
		return nil, fmt.Errorf("temporal: %w", err)
	}

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 20,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	orderRepo := orderPg.NewRepository(pgService)
	outboxRepo := outboxPg.NewRepository(pgService)
	txManager := pgx.NewTxManager(pgService.GetPool())

	return &InfraModule{
		TemporalClient:  temporalClient,
		PGService:       pgService,
		OrderRepo:       orderRepo,
		OutboxRepo:      outboxRepo,
		UserClient:      grpc_client.NewUserClient(httpClient, cfg.Upstream.UserCommandURL),
		ProductClient:   grpc_client.NewProductClient(httpClient, cfg.Upstream.ProductQueryURL),
		InventoryClient: grpc_client.NewInventoryClient(httpClient, cfg.Upstream.InventoryCommandURL),
		TxManager:       txManager,
		WorkflowRunner:  runner.NewWorkflowRunner(temporalClient, cfg.Temporal),
	}, nil
}
