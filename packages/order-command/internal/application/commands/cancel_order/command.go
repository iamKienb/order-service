package cancel_order

import (
	"context"

	"order-command-module/internal/domain/order"
	"order-command-module/internal/domain/shared"
)

type Command struct {
	OrderID   string
	ActorID   shared.UserID
	ActorType order.ActorType
	Reason    string
}

type Result struct {
	OrderID string
	Status  string
}

type Executor interface {
	Execute(ctx context.Context, cmd Command) (*Result, error)
}
