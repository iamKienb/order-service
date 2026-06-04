package port

import (
	"context"
	"order-command-module/internal/domain/shared"
)

type GetStockBySkuIDsRequest struct {
	ShopID string
	SkuIDs []string
}

type SkuStock struct {
	InventoryID       shared.InventoryID
	SkuID             shared.SkuID
	AvailableQuantity int64
}

type InventoryClient interface {
	GetStockBySkuIDs(ctx context.Context, req GetStockBySkuIDsRequest) ([]*SkuStock, error)
}
