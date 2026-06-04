package delete_inventories

import "context"

type inventoryService interface {
	DeleteInventories(ctx context.Context, cmd Command) (*Result, error)
}

type handler struct {
	service inventoryService
}

func NewHandler(service inventoryService) Executor {
	return &handler{service: service}
}

func (h *handler) Execute(ctx context.Context, cmd Command) (*Result, error) {
	return h.service.DeleteInventories(ctx, cmd)
}
