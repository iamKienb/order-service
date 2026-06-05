package workflow

import (
	"fmt"

	"order-command-module/internal/application/commands/place_order"
	"order-command-module/internal/bootstrap/config"
	domain_order "order-command-module/internal/domain/order"
	"order-command-module/internal/infra/temporal/activity"

	"go.temporal.io/sdk/workflow"
)

func PlaceOrderWorkflow(ctx workflow.Context, cmd place_order.Command, cfg config.TemporalConfig) (*place_order.Result, error) {
	activityCtx := activityContext(ctx, cfg)
	rollbackCtx := rollbackContext(ctx, cfg)

	var orderAct *activity.OrderActivity
	var result place_order.Result
	if err := workflow.ExecuteActivity(activityCtx, orderAct.CreateOrder, cmd).Get(ctx, &result); err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}
	if result.Status != string(domain_order.StatusPending) {
		return &result, nil
	}

	reserveCmd := activity.ReserveStockCommand{
		ShopID:  cmd.ShopID.String(),
		OrderID: result.OrderID,
		Items:   result.ReserveItems,
	}
	if err := workflow.ExecuteActivity(activityCtx, orderAct.ReserveStock, reserveCmd).Get(ctx, nil); err != nil {
		disconnectedCtx, _ := workflow.NewDisconnectedContext(rollbackCtx)
		_ = workflow.ExecuteActivity(disconnectedCtx, orderAct.CancelPlacedOrder, result.OrderID).Get(disconnectedCtx, nil)
		return nil, fmt.Errorf("reserve stock: %w", err)
	}

	return &result, nil
}
