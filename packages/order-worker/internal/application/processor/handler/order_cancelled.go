package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"order-shared-module/events"

	"order-worker-module/internal/application/port"
)

type OrderCancelledHandler struct {
	repo  port.ESRepository
	index string
}

func NewOrderCancelledHandler(repo port.ESRepository, index string) *OrderCancelledHandler {
	return &OrderCancelledHandler{repo: repo, index: index}
}

func (h *OrderCancelledHandler) Handle(ctx context.Context, raw json.RawMessage) error {
	var payload events.OrderCancelledEvent
	if err := json.Unmarshal(raw, &payload); err != nil {
		return fmt.Errorf("decode order cancelled event: %w", err)
	}
	if payload.OrderID == "" {
		return nil
	}

	doc := map[string]any{
		"id":       payload.OrderID,
		"shop_id":  payload.ShopID,
		"buyer_id": payload.BuyerID,

		"status":       payload.Status,
		"reason":       payload.Reason,
		"cancelled_at": payload.CancelledAt,
	}

	return h.repo.SyncData(ctx, h.index, payload.OrderID, doc)
}
