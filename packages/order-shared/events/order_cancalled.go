package events

const TopicOrderCancelled = "order-service.order.cancelled"

type OrderCancelledEvent struct {
	OrderID     string `json:"order_id"`
	ShopID      string `json:"shop_id"`
	BuyerID     string `json:"buyer_id"`
	Status      string `json:"status"`
	Reason      string `json:"reason"`
	CancelledAt string `json:"cancelled_at"`
}
