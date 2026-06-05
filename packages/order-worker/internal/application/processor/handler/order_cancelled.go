package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"order-worker-module/internal/application/port"
)

type OrderCancelledHandler struct {
	repo  port.ESRepository
	index string
}

type orderCancelledPatch struct {
	OrderID     string `json:"order_id"`
	ShopID      string `json:"shop_id"`
	BuyerID     string `json:"buyer_id"`
	Status      string `json:"status"`
	Reason      string `json:"reason"`
	CancelledAt string `json:"cancelled_at"`
}

func NewOrderCancelledHandler(repo port.ESRepository, index string) *OrderCancelledHandler {
	return &OrderCancelledHandler{repo: repo, index: index}
}

func (h *OrderCancelledHandler) Handle(ctx context.Context, payload json.RawMessage) error {
	var patch orderCancelledPatch
	if err := json.Unmarshal(payload, &patch); err != nil {
		return fmt.Errorf("decode order cancelled event: %w", err)
	}
	if patch.OrderID == "" {
		return nil
	}

	return h.repo.SyncData(ctx, h.index, patch.OrderID, patch)
}
