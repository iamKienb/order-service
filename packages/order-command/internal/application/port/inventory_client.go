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

type ReleaseAndFullfilStockParam struct {
	OrderID string
	ActorID string
}

type InventoryClient interface {
	GetStockBySkuIDs(ctx context.Context, req GetStockBySkuIDsRequest) ([]*SkuStock, error)
	ReserveStock(ctx context.Context, req ReserveStockRequest) error
	ReleaseStock(ctx context.Context, param ReleaseAndFullfilStockParam) error
	FulfillStock(ctx context.Context, param ReleaseAndFullfilStockParam) error
}

type ReserveStockRequest struct {
	ShopID  string
	OrderID string
	BuyerID string
	Items   []ReserveStockItem
}

type ReserveStockItem struct {
	InventoryID string
	SkuID       string
	Quantity    int64
}
