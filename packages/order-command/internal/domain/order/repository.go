package order

import (
	"context"

	"order-command-module/internal/domain/shared"
)

type QueryRepository interface {
	GetOrderByID(ctx context.Context, orderID string) (*Order, error)
	GetOrderByBuyerAndIdempotencyKey(ctx context.Context, buyerID shared.UserID, idempotencyKey string) (*Order, error)
}

type CommandRepository interface {
	CreateOrder(ctx context.Context, order *Order) error
	SaveOrderState(ctx context.Context, order *Order, expectedStatus OrderStatus) error
}

type Repository interface {
	QueryRepository
	CommandRepository
}
