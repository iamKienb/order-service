package order

import "context"

type QueryRepository interface {
	GetOrderByID(ctx context.Context, orderID string) (*Order, error)
}

type CommandRepository interface {
	CreateOrder(ctx context.Context, order *Order) error
	SaveOrderState(ctx context.Context, order *Order, expectedStatus OrderStatus) error
}

type Repository interface {
	QueryRepository
	CommandRepository
}
