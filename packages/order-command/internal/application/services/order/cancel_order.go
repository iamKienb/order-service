package order

import (
	"context"

	"order-command-module/internal/application/commands/cancel_order"
	"order-command-module/internal/application/port"
	domain_order "order-command-module/internal/domain/order"
)

func (s *orderService) CancelOrder(ctx context.Context, cmd cancel_order.Command) (*cancel_order.Result, error) {
	existing, err := s.orderRepo.GetOrderByID(ctx, cmd.OrderID)
	if err != nil {
		return nil, err
	}
	if existing.Status == domain_order.StatusCancelled {
		return &cancel_order.Result{OrderID: existing.ID, Status: string(domain_order.StatusCancelled)}, nil
	}

	expectedStatus := existing.Status
	actorID := cmd.ActorID
	if err := existing.Cancel(domain_order.CancelParams{
		ActorID:   &actorID,
		ActorType: string(cmd.ActorType),
		Reason:    cmd.Reason,
	}); err != nil {
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

	return &cancel_order.Result{OrderID: existing.ID, Status: string(domain_order.StatusCancelled)}, nil
}
