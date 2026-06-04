package inventory

import (
	"context"
	"fmt"
	domain_inventory "inventory-command-module/internal/domain/inventory"
	"inventory-command-module/internal/domain/shared"
)

func (r *inventoryRepository) GetListInventoryItemsBySkuIDs(ctx context.Context, skuIDs []shared.SkuID) ([]*domain_inventory.InventoryItem, error) {
	rows, err := r.getQuerier(ctx).ListInventoryItemsBySkuIDs(ctx, toUUUIDs(skuIDs))
	if err != nil {
		return nil, fmt.Errorf("infra: list inventory items by sku ids: %w", err)
	}

	return toDomainInventories(rows), nil
}
