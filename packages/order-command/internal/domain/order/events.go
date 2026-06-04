package order

import (
	"time"

	"order-command-module/internal/domain/shared"
)

type InventoryCreatedEvent struct {
	InventoryID shared.InventoryID
	SkuID       shared.SkuID
	ShopID      shared.ShopID
	Quantity    int64
	Reserved    int64
	Status      OrderStatus

	CreatedBy shared.UserID
	CreatedAt time.Time
}

func (e InventoryCreatedEvent) EventName() string {
	return "inventory-service.inventory.created"
}

func (e InventoryCreatedEvent) IntegrationPayload() map[string]interface{} {
	return map[string]interface{}{
		"inventory_id": e.InventoryID.String(),
		"sku_id":       e.SkuID.String(),
		"shop_id":      e.ShopID.String(),
		"quantity":     e.Quantity,
		"reserved":     e.Reserved,
		"status":       string(e.Status),
		"created_by":   e.CreatedBy.String(),
		"created_at":   e.CreatedAt,
	}
}
