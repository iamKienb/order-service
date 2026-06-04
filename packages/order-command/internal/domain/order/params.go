package order

import (
	"order-command-module/internal/domain/shared"
	"time"
)

type NewOrderParams struct {
	ShopID           shared.ShopID
	BuyerID          shared.UserID
	ShippingName     string
	ShippingPhone    string
	ShippingAddress  string
	ShippingProvince string
	ShippingWard     string
	Note             *string
	GrandTotal       int64
	Currency         string
}

type OrderItemsParams struct {
	ProductID   shared.ProductID
	InventoryID shared.InventoryID
	SkuID       shared.SkuID

	ProductName string
	SkuCode     string
	ImageUrl    string
	Quantity    int64

	BasePrice    int64
	ItemSubtotal int64
}

type NewOrderItemParams struct {
	OrderID     string
	OrderItemID shared.OrderItemID
	ProductID   shared.ProductID
	InventoryID shared.InventoryID
	SkuID       shared.SkuID

	ProductName string
	SkuCode     string
	ImageUrl    string
	Quantity    int64

	BasePrice int64

	Currency  string
	CreatedAt time.Time
}

type NewOrderHistoryParams struct {
	ID         shared.OrderHistoryID
	OrderID    string
	FromStatus *OrderStatus
	ToStatus   OrderStatus
	ChangedBy  shared.UserID
	ActorType  ActorType
	Reason     string
	CreatedAt  time.Time
}

type CancelParams struct {
	ActorID   *shared.UserID
	Reason    string
	ActorType string
}
