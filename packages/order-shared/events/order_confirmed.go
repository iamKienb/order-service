package events

const TopicOrderConfirmed = "order-service.order.confirmed"

type OrderConfirmedEvent struct {
	OrderID     string `json:"order_id"`
	ShopID      string `json:"shop_id"`
	Status      string `json:"status"`
	ConfirmedAt string `json:"confirmed_at"`
}
