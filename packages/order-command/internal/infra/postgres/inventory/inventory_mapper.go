package inventory

import (
	"inventory-command-module/db/repository"
	domain_inventory "inventory-command-module/internal/domain/inventory"
	"inventory-command-module/internal/domain/shared"

	"github.com/iamKienb/go-core/postgres/conv"
	"github.com/jackc/pgx/v5/pgtype"
)

func toDomainInventories(rows []repository.Inventory) []*domain_inventory.InventoryItem {
	inventories := make([]*domain_inventory.InventoryItem, 0, len(rows))
	for _, row := range rows {
		inventories = append(inventories, toDomainInventory(row))
	}

	return inventories
}

func toDomainInventory(row repository.Inventory) *domain_inventory.InventoryItem {
	return &domain_inventory.InventoryItem{
		ID:        shared.InventoryID(row.ID.Bytes),
		SkuID:     shared.SkuID(row.SkuID.Bytes),
		ShopID:    shared.ShopID(row.ShopID.Bytes),
		Quantity:  row.Quantity,
		Reserved:  row.Reserved,
		Status:    domain_inventory.InventoryStatus(row.Status),
		CreatedBy: shared.UserID(row.CreatedBy.Bytes),
		UpdatedBy: toUserIDPointer(row.UpdatedBy),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: toTimePointer(row.UpdatedAt),
	}
}

func toInfraCreateInventoryBatch(inventories []*domain_inventory.InventoryItem) repository.CreateInventoryItemsBatchParams {
	params := repository.CreateInventoryItemsBatchParams{
		Ids:                make([]pgtype.UUID, 0, len(inventories)),
		SkuIds:             make([]pgtype.UUID, 0, len(inventories)),
		Quantities:         make([]int64, 0, len(inventories)),
		ReservedQuantities: make([]int64, 0, len(inventories)),
		Statuses:           make([]string, 0, len(inventories)),
		CreatedAt:          conv.TimeStampZ(createdAtOf(inventories)),
	}

	for _, inventory := range inventories {
		params.Ids = append(params.Ids, conv.UUID(inventory.ID))
		params.SkuIds = append(params.SkuIds, conv.UUID(inventory.SkuID))
		params.Quantities = append(params.Quantities, inventory.Quantity)
		params.ReservedQuantities = append(params.ReservedQuantities, inventory.Reserved)
		params.Statuses = append(params.Statuses, string(inventory.Status))
	}

	if len(inventories) > 0 {
		params.ShopID = conv.UUID(inventories[0].ShopID)
		params.CreatedBy = conv.UUID(inventories[0].CreatedBy)
		params.UpdatedBy = conv.UUID(inventories[0].CreatedBy)
	}

	return params
}

func toInfraReserveInventoryStockBatch(params domain_inventory.ReserveStockParams) repository.ReserveInventoryStockBatchParams {
	return repository.ReserveInventoryStockBatchParams{
		ShopID:       conv.UUID(params.ShopID),
		ActorID:      conv.UUID(params.ActorID),
		InventoryIds: toUUUIDs(params.InventoryIds),
		Quantities:   params.Quantities,
		CreatedAt:    conv.TimeStampZ(&params.CreatedAt),
	}
}

func toInfraReleaseInventoryStockBatch(params domain_inventory.ReleaseAndFulfillStockParams) repository.ReleaseInventoryStockBatchParams {
	return repository.ReleaseInventoryStockBatchParams{
		ShopID:       conv.UUID(params.ShopID),
		ActorID:      conv.UUID(params.ActorID),
		InventoryIds: toUUUIDs(params.InventoryIds),
		Quantities:   params.Quantities,
		UpdatedAt:    conv.TimeStampZ(&params.UpdatedAt),
	}
}

func toInfraFulfillInventoryStockBatch(params domain_inventory.ReleaseAndFulfillStockParams) repository.FulfillInventoryStockBatchParams {
	return repository.FulfillInventoryStockBatchParams{
		ShopID:       conv.UUID(params.ShopID),
		ActorID:      conv.UUID(params.ActorID),
		InventoryIds: toUUUIDs(params.InventoryIds),
		Quantities:   params.Quantities,
		UpdatedAt:    conv.TimeStampZ(&params.UpdatedAt),
	}
}

func toInfraCreateTransactionBatch(transactions []*domain_inventory.InventoryTransaction) repository.CreateInventoryTransactionBatchParams {
	params := repository.CreateInventoryTransactionBatchParams{
		Ids:          make([]pgtype.UUID, 0, len(transactions)),
		InventoryIds: make([]pgtype.UUID, 0, len(transactions)),

		Types:          make([]string, len(transactions)),
		Quantities:     make([]int64, len(transactions)),
		BalancesBefore: make([]int64, len(transactions)),
		BalancesAfter:  make([]int64, len(transactions)),

		ReferenceTypes: make([]string, len(transactions)),
		ReferenceIds:   make([]string, len(transactions)),
		ActionTypes:    make([]string, len(transactions)),

		IdempotencyKeys: make([]string, len(transactions)),
		Notes:           make([]string, len(transactions)),
		CreatedBys:      make([]pgtype.UUID, len(transactions)),
		CreatedAts:      make([]pgtype.Timestamptz, len(transactions)),
	}

	for _, transaction := range transactions {
		params.Ids = append(params.Ids, conv.UUID(transaction.ID))
		params.InventoryIds = append(params.InventoryIds, conv.UUID(transaction.InventoryID))

		params.Types = append(params.Types, string(transaction.Type))
		params.Quantities = append(params.Quantities, transaction.Quantity)
		params.BalancesBefore = append(params.BalancesBefore, transaction.BalanceBefore)
		params.BalancesAfter = append(params.BalancesAfter, transaction.BalanceBefore)

		params.ReferenceTypes = append(params.ReferenceTypes, transaction.ReferenceType)
		params.ReferenceIds = append(params.ReferenceIds, transaction.ReferenceID)
		params.ActionTypes = append(params.ActionTypes, string(transaction.ActionType))

		params.IdempotencyKeys = append(params.IdempotencyKeys, transaction.IdempotencyKey)
		params.Notes = append(params.Notes, transaction.Note)
		params.CreatedBys = append(params.CreatedBys, conv.UUID(transaction.CreatedBy))
		params.CreatedAts = append(params.CreatedAts, conv.TimeStampZ(&transaction.CreatedAt))
	}

	return params
}
