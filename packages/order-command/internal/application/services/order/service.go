package order

import (
	"context"
	"order-command-module/internal/application/commands/preview_checkout"
)

type Service interface {
	PreviewCheckout(ctx context.Context, cmd preview_checkout.Command, checkoutCtx preview_checkout.CheckoutContext) (*preview_checkout.Result, error)
}

type orderService struct {
	// orderRepo

}

func NewOrderService() Service {
	return &orderService{}
}
