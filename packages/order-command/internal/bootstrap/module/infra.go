package module

import (
	"context"
	"fmt"

	"inventory-command-module/internal/application/port"
	"inventory-command-module/internal/bootstrap/config"
	"inventory-command-module/internal/domain/inventory"
	inventoryPg "inventory-command-module/internal/infra/postgres/inventory"
	outboxPg "inventory-command-module/internal/infra/postgres/outbox"

	pgx "github.com/iamKienb/go-core/postgres"
)

type InfraModule struct {
	PGService pgx.PGXService

	InventoryRepo inventory.Repository
	OutboxRepo    port.OutboxRepository
	TxManager     port.TxManager
}

func NewInfraModule(ctx context.Context, cfg *config.InventoryCommandConfig) (*InfraModule, error) {
	pgService, err := pgx.New(cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}

	return &InfraModule{
		PGService:     pgService,
		InventoryRepo: inventoryPg.NewRepository(pgService),
		OutboxRepo:    outboxPg.NewRepository(pgService),
		TxManager:     pgx.NewTxManager(pgService.GetPool()),
	}, nil
}
