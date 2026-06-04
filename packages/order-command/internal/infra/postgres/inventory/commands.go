package inventory

import (
	"context"
	"fmt"
	"time"

	"inventory-command-module/db/repository"
	domain_inventory "inventory-command-module/internal/domain/inventory"
	"inventory-command-module/internal/domain/shared"

	"github.com/iamKienb/go-core/postgres/conv"
)

func (r *inventoryRepository) CreateInventoryItemsBatch(ctx context.Context, inventories []*domain_inventory.InventoryItem) error {
	if len(inventories) == 0 {
		return nil
	}
	if err := r.getQuerier(ctx).CreateInventoryItemsBatch(ctx, toInfraCreateInventoryBatch(inventories)); err != nil {
		if r.isDuplicateSKU(err) {
			return domain_inventory.ErrInventoryItemAlreadyExists
		}
		return fmt.Errorf("infra: create inventory items batch: %w", err)
	}
	return nil
}

func (r *inventoryRepository) CreateTransactionBatch(ctx context.Context, transactions []*domain_inventory.InventoryTransaction) error {
	if len(transactions) == 0 {
		return nil
	}
	if err := r.getQuerier(ctx).CreateInventoryTransactionBatch(ctx, toInfraCreateTransactionBatch(transactions)); err != nil {
		return fmt.Errorf("infra: save inventory transaction batch: %w", err)
	}
	return nil
}

func (r *inventoryRepository) SoftDeleteBySkuIDs(ctx context.Context, skuIDs []shared.SkuID, actorID shared.UserID, deletedAt time.Time) error {
	rowsAffected, err := r.getQuerier(ctx).SoftDeleteInventoryItemsBySkuIDs(ctx, repository.SoftDeleteInventoryItemsBySkuIDsParams{
		DeletedStatus: string(domain_inventory.StatusDeleted),
		SkuIds:        toUUUIDs(skuIDs),
		ActorID:       conv.UUID(actorID),
		UpdatedAt:     conv.TimeStampZ(&deletedAt),
		ActiveStatus:  string(domain_inventory.StatusActive),
	})
	if err != nil {
		return fmt.Errorf("infra: soft delete inventory items by sku ids: %w", err)
	}
	if rowsAffected == 0 {
		return domain_inventory.ErrInventoryItemNotFound
	}

	return nil
}

func (r *inventoryRepository) ReserveInventoryStockBatch(ctx context.Context, params domain_inventory.ReserveStockParams) (int64, error) {
	RowsAffected, err := r.getQuerier(ctx).ReserveInventoryStockBatch(ctx, toInfraReserveInventoryStockBatch(params))
	if err != nil {
		return RowsAffected, fmt.Errorf("infra: reserve inventory item: %w", err)
	}

	return RowsAffected, nil
}

func (r *inventoryRepository) ReleaseInventoryStockBatch(ctx context.Context, params domain_inventory.ReleaseAndFulfillStockParams) (int64, error) {
	releaseRows, err := r.getQuerier(ctx).ReleaseInventoryStockBatch(ctx, toInfraReleaseInventoryStockBatch(params))
	if err != nil {
		return releaseRows, fmt.Errorf("infra: release reservations by order id: %w", err)
	}

	return releaseRows, nil
}

func (r *inventoryRepository) FulfillInventoryStockBatch(ctx context.Context, params domain_inventory.ReleaseAndFulfillStockParams) (int64, error) {
	releaseRows, err := r.getQuerier(ctx).FulfillInventoryStockBatch(ctx, toInfraFulfillInventoryStockBatch(params))
	if err != nil {
		return releaseRows, fmt.Errorf("infra: fullfil reservations by order id: %w", err)
	}

	return releaseRows, nil
}

func (r *inventoryRepository) CreateReservationBatch(ctx context.Context, reservations []*domain_inventory.InventoryReservation) error {
	if len(reservations) == 0 {
		return nil
	}

	if err := r.getQuerier(ctx).CreateInventoryReservationBatch(ctx, toInfraCreateReservationBatch(reservations)); err != nil {
		return fmt.Errorf("infra: save inventory reservation: %w", err)
	}

	return nil
}

func (r *inventoryRepository) ReleaseReservationsByOrderID(ctx context.Context, params domain_inventory.ReleaseAndFulfillReservationParams) ([]*domain_inventory.InventoryReservation, error) {
	reservations, err := r.getQuerier(ctx).ReleaseReservationsByOrderID(ctx, toInfraReleaseReservationsByOrderID(params))
	if err != nil {
		return nil, fmt.Errorf("infra: release reservations by order id: %w", err)
	}

	return toDomainInventoryReservations(reservations), nil
}

func (r *inventoryRepository) FulfillReservationsByOrderID(ctx context.Context, params domain_inventory.ReleaseAndFulfillReservationParams) ([]*domain_inventory.InventoryReservation, error) {
	reservations, err := r.getQuerier(ctx).FulfillReservationsByOrderID(ctx, toInfraFulfillReservationsByOrderID(params))
	if err != nil {
		return nil, fmt.Errorf("infra: FullFill reservations by order id: %w", err)
	}

	return toDomainInventoryReservations(reservations), nil
}
