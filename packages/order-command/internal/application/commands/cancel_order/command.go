package create_inventories

import (
	"context"

	"inventory-command-module/internal/domain/shared"
)

type Item struct {
	SkuID    shared.SkuID
	Quantity int64
}

type Command struct {
	ShopID    shared.ShopID
	ActorID   shared.UserID
	Inventory []Item
}

type Result struct {
	Success bool
}

type Executor interface {
	Execute(ctx context.Context, cmd Command) (*Result, error)
}
