package confirm_order

import (
	"context"

	"order-command-module/internal/domain/shared"
)

type Command struct {
	OrderID string
	ShopID  shared.ShopID
	ActorID shared.UserID
}

type Result struct {
	OrderID string
	Status  string
}

type Executor interface {
	Execute(ctx context.Context, cmd Command) (*Result, error)
}
