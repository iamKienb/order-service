package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"order-worker-module/internal/application/port"
)

type OrderConfirmedHandler struct {
	repo  port.ESRepository
	index string
}

type orderConfirmedPatch struct {
	OrderID     string `json:"order_id"`
	ShopID      string `json:"shop_id"`
	Status      string `json:"status"`
	ConfirmedAt string `json:"confirmed_at"`
}

func NewOrderConfirmedHandler(repo port.ESRepository, index string) *OrderConfirmedHandler {
	return &OrderConfirmedHandler{repo: repo, index: index}
}

func (h *OrderConfirmedHandler) Handle(ctx context.Context, payload json.RawMessage) error {
	var patch orderConfirmedPatch
	if err := json.Unmarshal(payload, &patch); err != nil {
		return fmt.Errorf("decode order confirmed event: %w", err)
	}
	if patch.OrderID == "" {
		return nil
	}

	return h.repo.SyncData(ctx, h.index, patch.OrderID, patch)
}
