package order

import (
	"order-command-module/internal/domain/shared"
	"time"
)

type OrderItem struct {
	ID          shared.OrderItemID
	OrderID     string
	ProductID   shared.ProductID
	InventoryID shared.InventoryID
	SkuID       shared.SkuID

	ProductName string
	SkuCode     string
	ImageUrl    string
	Quantity    int64

	BasePrice    int64
	ItemSubtotal int64
	Currency     string
	CreatedAt    time.Time
}

func createOrderItem(params NewOrderItemParams) *OrderItem {
	itemSubtotal := params.BasePrice * params.Quantity

	return &OrderItem{
		ID:          params.OrderItemID,
		OrderID:     params.OrderID,
		ProductID:   params.ProductID,
		InventoryID: params.InventoryID,
		SkuID:       params.SkuID,
		ProductName: params.ProductName,
		SkuCode:     params.SkuCode,
		ImageUrl:    params.ImageUrl,
		Quantity:    params.Quantity,

		BasePrice:    params.BasePrice,
		ItemSubtotal: itemSubtotal,
		Currency:     params.Currency,
		CreatedAt:    params.CreatedAt,
	}
}
