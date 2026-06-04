package create_inventories

import "context"

type inventoryService interface {
	CreateInventories(ctx context.Context, cmd Command) (*Result, error)
}

type handler struct {
	service inventoryService
}

func NewHandler(service inventoryService) Executor {
	return &handler{service: service}
}

func (h *handler) Execute(ctx context.Context, cmd Command) (*Result, error) {
	return h.service.CreateInventories(ctx, cmd)
}
