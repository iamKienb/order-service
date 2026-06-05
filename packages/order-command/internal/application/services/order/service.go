package order

import (
	"context"

	"order-command-module/internal/application/commands/cancel_order"
	"order-command-module/internal/application/commands/confirm_order"
	"order-command-module/internal/application/commands/place_order"
	"order-command-module/internal/application/commands/preview_checkout"
	"order-command-module/internal/application/port"
	"order-command-module/internal/application/services/outbox"
	domain_order "order-command-module/internal/domain/order"
)

type Service interface {
	PreviewCheckout(ctx context.Context, cmd preview_checkout.Command, checkoutCtx preview_checkout.CheckoutContext) (*preview_checkout.Result, error)
	PlaceOrder(ctx context.Context, cmd place_order.Command, checkoutCtx preview_checkout.CheckoutContext) (*place_order.Result, error)
	ConfirmOrder(ctx context.Context, cmd confirm_order.Command) (*confirm_order.Result, error)
	CancelOrder(ctx context.Context, cmd cancel_order.Command) (*cancel_order.Result, error)
}

type orderService struct {
	orderRepo     domain_order.Repository
	outboxService outbox.Service
	txManager     port.TxManager
}

func NewOrderService(orderRepo domain_order.Repository, outboxService outbox.Service, txManager port.TxManager) Service {
	return &orderService{
		orderRepo:     orderRepo,
		outboxService: outboxService,
		txManager:     txManager,
	}
}
