package search_orders

import "context"

type orderQueryService interface {
	SearchOrders(context.Context, Query) (*Result, error)
}

type handler struct {
	service orderQueryService
}

func NewHandler(service orderQueryService) Executor {
	return &handler{service: service}
}

func (h *handler) Execute(ctx context.Context, query Query) (*Result, error) {
	return h.service.SearchOrders(ctx, query)
}
