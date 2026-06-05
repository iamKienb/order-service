package place_order

import (
	"context"

	"order-command-module/internal/domain/shared"
)

type Item struct {
	SkuID     shared.SkuID
	BasePrice int64
	Quantity  int64
}

type Command struct {
	ShopID         shared.ShopID
	BuyerID        shared.UserID
	BuyerAddressID shared.UserAddressID
	IdempotencyKey string
	Items          []Item
}

type Result struct {
	OrderID      string
	Status       string
	ReserveItems []ReserveItem
}

type ReserveItem struct {
	InventoryID string
	SkuID       string
	Quantity    int64
}

type Executor interface {
	Execute(ctx context.Context, cmd Command) (*Result, error)
}
