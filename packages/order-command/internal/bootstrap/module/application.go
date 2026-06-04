package module

import (
	"inventory-command-module/internal/application/commands/create_inventories"
	"inventory-command-module/internal/application/commands/delete_inventories"
	"inventory-command-module/internal/application/commands/fulfill_stock"
	"inventory-command-module/internal/application/commands/release_stock"
	"inventory-command-module/internal/application/commands/reserve_stock"
	"inventory-command-module/internal/application/services/inventory"
	"inventory-command-module/internal/application/services/outbox"
)

type ApplicationModule struct {
	CreateInventoriesExecutor create_inventories.Executor
	DeleteInventoriesExecutor delete_inventories.Executor
	ReserveStockExecutor      reserve_stock.Executor
	ReleaseStockExecutor      release_stock.Executor
	FulfillStockExecutor      fulfill_stock.Executor
}

func NewApplicationModule(infra *InfraModule) *ApplicationModule {
	outboxService := outbox.NewOutboxService(infra.OutboxRepo)
	inventoryService := inventory.NewInventoryService(
		infra.InventoryRepo,
		outboxService,
		infra.TxManager,
	)

	return &ApplicationModule{
		CreateInventoriesExecutor: create_inventories.NewHandler(inventoryService),
		DeleteInventoriesExecutor: delete_inventories.NewHandler(inventoryService),
		ReserveStockExecutor:      reserve_stock.NewHandler(inventoryService),
		ReleaseStockExecutor:      release_stock.NewHandler(inventoryService),
		FulfillStockExecutor:      fulfill_stock.NewHandler(inventoryService),
	}
}
