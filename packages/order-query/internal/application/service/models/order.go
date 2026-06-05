package models

type Page struct {
	Size  int
	Token string
}

type Order struct {
	OrderID     string      `json:"order_id"`
	ShopID      string      `json:"shop_id"`
	BuyerID     string      `json:"buyer_id"`
	Status      string      `json:"status"`
	GrandTotal  int64       `json:"grand_total"`
	Currency    string      `json:"currency"`
	Reason      string      `json:"reason"`
	CreatedAt   string      `json:"created_at"`
	ConfirmedAt string      `json:"confirmed_at"`
	CancelledAt string      `json:"cancelled_at"`
	Items       []OrderItem `json:"items"`
}

type OrderItem struct {
	InventoryID  string `json:"inventory_id"`
	SkuID        string `json:"sku_id"`
	SkuCode      string `json:"sku_code"`
	ProductID    string `json:"product_id"`
	ProductName  string `json:"product_name"`
	Quantity     int64  `json:"quantity"`
	BasePrice    int64  `json:"base_price"`
	ItemSubtotal int64  `json:"item_subtotal"`
}

type OrderPage struct {
	Items         []Order
	Total         int64
	NextPageToken string
}
