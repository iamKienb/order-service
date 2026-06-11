package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"order-shared-module/events"

	"order-worker-module/internal/application/port"
)

type OrderConfirmedHandler struct {
	repo  port.ESRepository
	index string
}

func NewOrderConfirmedHandler(repo port.ESRepository, index string) *OrderConfirmedHandler {
	return &OrderConfirmedHandler{repo: repo, index: index}
}

func (h *OrderConfirmedHandler) Handle(ctx context.Context, raw json.RawMessage) error {
	var payload events.OrderConfirmedEvent
	if err := json.Unmarshal(raw, &payload); err != nil {
		return fmt.Errorf("decode order confirmed event: %w", err)
	}
	if payload.OrderID == "" {
		return nil
	}

	doc := map[string]any{
		"id":           payload.OrderID,
		"shop_id":      payload.ShopID,
		"status":       payload.Status,
		"confirmed_at": payload.ConfirmedAt,
	}

	return h.repo.SyncData(ctx, h.index, payload.OrderID, doc)
}
