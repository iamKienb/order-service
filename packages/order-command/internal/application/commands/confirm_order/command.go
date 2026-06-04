package delete_inventories

import (
	"context"

	"inventory-command-module/internal/domain/shared"
)

type Command struct {
	ActorID shared.UserID
	SkuIDs  []shared.SkuID
}

type Result struct {
	Success bool
}

type Executor interface {
	Execute(ctx context.Context, cmd Command) (*Result, error)
}
