package order

import (
	"context"

	"order-command-module/internal/application/commands/confirm_order"
	"order-command-module/internal/application/port"
	domain_order "order-command-module/internal/domain/order"
)

func (s *orderService) ConfirmOrder(ctx context.Context, cmd confirm_order.Command) (*confirm_order.Result, error) {
	existing, err := s.orderRepo.GetOrderByID(ctx, cmd.OrderID)
	if err != nil {
		return nil, err
	}
	if existing.ShopID != cmd.ShopID {
		return nil, domain_order.ErrOrderShopMismatch
	}
	if existing.Status == domain_order.StatusConfirmed {
		return &confirm_order.Result{OrderID: existing.ID, Status: string(domain_order.StatusConfirmed)}, nil
	}

	expectedStatus := existing.Status
	if err := existing.Confirm(cmd.ActorID); err != nil {
		return nil, err
	}

	if err := s.txManager.WithTx(ctx, func(ctx context.Context) error {
		if err := s.orderRepo.SaveOrderState(ctx, existing, expectedStatus); err != nil {
			return err
		}
		if events := existing.FlushEvents(); len(events) > 0 {
			return s.outboxService.PublishBatch(ctx, []port.OutboxParam{{
				AggregateID:   orderAggregateID(existing.ID),
				AggregateType: existing.Type(),
				Events:        events,
			}})
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &confirm_order.Result{OrderID: existing.ID, Status: string(domain_order.StatusConfirmed)}, nil
}
