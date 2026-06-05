package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"order-worker-module/internal/application/port"
)

type OrderCreatedHandler struct {
	repo  port.ESRepository
	index string
}

type orderCreatedDocument struct {
	OrderID    string                     `json:"order_id"`
	ShopID     string                     `json:"shop_id"`
	BuyerID    string                     `json:"buyer_id"`
	Status     string                     `json:"status"`
	GrandTotal int64                      `json:"grand_total"`
	Currency   string                     `json:"currency"`
	Items      []orderCreatedDocumentItem `json:"items"`
	CreatedAt  string                     `json:"created_at"`
}

type orderCreatedDocumentItem struct {
	InventoryID  string `json:"inventory_id"`
	SkuID        string `json:"sku_id"`
	SkuCode      string `json:"sku_code"`
	ProductID    string `json:"product_id"`
	ProductName  string `json:"product_name"`
	Quantity     int64  `json:"quantity"`
	BasePrice    int64  `json:"base_price"`
	ItemSubtotal int64  `json:"item_subtotal"`
}

func NewOrderCreatedHandler(repo port.ESRepository, index string) *OrderCreatedHandler {
	return &OrderCreatedHandler{repo: repo, index: index}
}

func (h *OrderCreatedHandler) Handle(ctx context.Context, payload json.RawMessage) error {
	var doc orderCreatedDocument
	if err := json.Unmarshal(payload, &doc); err != nil {
		return fmt.Errorf("decode order created event: %w", err)
	}
	if doc.OrderID == "" {
		return nil
	}

	return h.repo.SyncData(ctx, h.index, doc.OrderID, doc)
}
