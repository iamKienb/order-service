package order

import "time"

type OrderCreatedEvent struct {
	OrderID    string
	ShopID     string
	BuyerID    string
	Status     OrderStatus
	GrandTotal int64
	Currency   string
	Items      []OrderCreatedEventItem
	CreatedAt  time.Time
}

type OrderCreatedEventItem struct {
	InventoryID  string
	SkuID        string
	SkuCode      string
	ProductID    string
	ProductName  string
	Quantity     int64
	BasePrice    int64
	ItemSubtotal int64
}

func (e OrderCreatedEvent) EventName() string {
	return "order-service.order.created"
}

func (e OrderCreatedEvent) IntegrationPayload() map[string]interface{} {
	items := make([]map[string]interface{}, 0, len(e.Items))
	for _, item := range e.Items {
		items = append(items, map[string]interface{}{
			"inventory_id":  item.InventoryID,
			"sku_id":        item.SkuID,
			"sku_code":      item.SkuCode,
			"product_id":    item.ProductID,
			"product_name":  item.ProductName,
			"quantity":      item.Quantity,
			"base_price":    item.BasePrice,
			"item_subtotal": item.ItemSubtotal,
		})
	}

	return map[string]interface{}{
		"order_id":    e.OrderID,
		"shop_id":     e.ShopID,
		"buyer_id":    e.BuyerID,
		"status":      string(e.Status),
		"grand_total": e.GrandTotal,
		"currency":    e.Currency,
		"items":       items,
		"created_at":  e.CreatedAt,
	}
}

type OrderConfirmedEvent struct {
	OrderID     string
	ShopID      string
	Status      OrderStatus
	ConfirmedAt time.Time
}

func (e OrderConfirmedEvent) EventName() string {
	return "order-service.order.confirmed"
}

func (e OrderConfirmedEvent) IntegrationPayload() map[string]interface{} {
	return map[string]interface{}{
		"order_id":     e.OrderID,
		"shop_id":      e.ShopID,
		"status":       string(e.Status),
		"confirmed_at": e.ConfirmedAt,
	}
}

type OrderCancelledEvent struct {
	OrderID     string
	ShopID      string
	BuyerID     string
	Status      OrderStatus
	Reason      string
	CancelledAt time.Time
}

func (e OrderCancelledEvent) EventName() string {
	return "order-service.order.cancelled"
}

func (e OrderCancelledEvent) IntegrationPayload() map[string]interface{} {
	return map[string]interface{}{
		"order_id":     e.OrderID,
		"shop_id":      e.ShopID,
		"buyer_id":     e.BuyerID,
		"status":       string(e.Status),
		"reason":       e.Reason,
		"cancelled_at": e.CancelledAt,
	}
}
