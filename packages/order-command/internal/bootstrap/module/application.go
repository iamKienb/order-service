package module

import (
	"order-command-module/internal/application/commands/cancel_order"
	"order-command-module/internal/application/commands/confirm_order"
	"order-command-module/internal/application/commands/place_order"
	"order-command-module/internal/application/commands/preview_checkout"
	orderapp "order-command-module/internal/application/services/order"
	"order-command-module/internal/application/services/outbox"
)

type ApplicationModule struct {
	OrderService            orderapp.Service
	PreviewCheckoutExecutor preview_checkout.Executor
	PlaceOrderExecutor      place_order.Executor
	ConfirmOrderExecutor    confirm_order.Executor
	CancelOrderExecutor     cancel_order.Executor
}

func NewApplicationModule(infra *InfraModule) *ApplicationModule {
	outboxService := outbox.NewOutboxService(infra.OutboxRepo)
	orderService := orderapp.NewOrderService(
		infra.OrderRepo,
		outboxService,
		infra.TxManager,
	)

	return &ApplicationModule{
		OrderService: orderService,
		PreviewCheckoutExecutor: preview_checkout.NewHandler(
			orderService,
			infra.UserClient,
			infra.ProductClient,
			infra.InventoryClient,
		),
		PlaceOrderExecutor:   place_order.NewHandler(infra.WorkflowRunner),
		ConfirmOrderExecutor: confirm_order.NewHandler(infra.WorkflowRunner),
		CancelOrderExecutor:  cancel_order.NewHandler(infra.WorkflowRunner),
	}
}
