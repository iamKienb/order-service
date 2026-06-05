package preview_checkout

import (
	"context"

	"order-command-module/internal/domain/shared"
)

type Item struct {
	SkuID    shared.SkuID
	Quantity int64
}

type Command struct {
	ShopID         shared.ShopID
	BuyerID        shared.UserID
	BuyerAddressID shared.UserAddressID
	Items          []Item
}

type PreviewItemDetail struct {
	ShopID      shared.ShopID
	InventoryID shared.InventoryID
	SkuID       shared.SkuID
	SkuCode     string
	ProductID   shared.ProductID
	ProductName string
	BasePrice   int64
	SubTotal    int64
	Quantity    int64
	ImageURL    string
}

type Result struct {
	GrandTotal int64
	Items      []PreviewItemDetail
}

type Executor interface {
	Execute(ctx context.Context, cmd Command) (*Result, error)
}
