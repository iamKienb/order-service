package list_buyer_orders

import "context"

type orderQueryService interface {
	ListBuyerOrders(context.Context, Query) (*Result, error)
}

type handler struct {
	service orderQueryService
}

func NewHandler(service orderQueryService) Executor {
	return &handler{service: service}
}

func (h *handler) Execute(ctx context.Context, query Query) (*Result, error) {
	return h.service.ListBuyerOrders(ctx, query)
}
