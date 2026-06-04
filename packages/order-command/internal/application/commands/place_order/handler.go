package place_order

import "context"

type CheckoutContext struct {
	ProductSkus []*port.ProductSku
	SkuStocks   []*port.SkuStock
	UserAddress *port.UserAddress
}

type orderService interface {
	PlaceOrder(ctx context.Context, cmd Command) (*Result, error)
}

type handler struct {
	service orderService
}

func NewHandler(service orderService) Executor {
	return &handler{service: service}
}

func (h *handler) Execute(ctx context.Context, cmd Command) (*Result, error) {
	return h.service.PlaceOrder(ctx, cmd)
}
