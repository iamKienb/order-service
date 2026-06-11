package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"order-shared-module/events"

	"order-worker-module/internal/application/port"
)

type OrderCreatedHandler struct {
	repo  port.ESRepository
	index string
}

func NewOrderCreatedHandler(repo port.ESRepository, index string) *OrderCreatedHandler {
	return &OrderCreatedHandler{repo: repo, index: index}
}

func (h *OrderCreatedHandler) Handle(ctx context.Context, raw json.RawMessage) error {
	var payload events.OrderCreatedEvent
	if err := json.Unmarshal(raw, &payload); err != nil {
		return fmt.Errorf("decode order created event: %w", err)
	}
	if payload.OrderID == "" {
		return nil
	}

	doc := map[string]any{
		"id":          payload.OrderID,
		"shop_id":     payload.ShopID,
		"buyer_id":    payload.BuyerID,
		"status":      payload.Status,
		"grand_total": payload.GrandTotal,
		"currency":    payload.Currency,
		"items":       payload.Items,
		"created_at":  payload.CreatedAt,
	}

	return h.repo.SyncData(ctx, h.index, payload.OrderID, doc)
}
