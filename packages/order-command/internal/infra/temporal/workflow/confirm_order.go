package workflow

import (
	"fmt"

	"order-command-module/internal/application/commands/confirm_order"
	"order-command-module/internal/bootstrap/config"
	"order-command-module/internal/infra/temporal/activity"

	"go.temporal.io/sdk/workflow"
)

func ConfirmOrderWorkflow(ctx workflow.Context, cmd confirm_order.Command, cfg config.TemporalConfig) (*confirm_order.Result, error) {
	activityCtx := activityContext(ctx, cfg)
	rollbackCtx := rollbackContext(ctx, cfg)

	var orderAct *activity.OrderActivity
	var result confirm_order.Result
	if err := workflow.ExecuteActivity(activityCtx, orderAct.ConfirmOrder, cmd).Get(ctx, &result); err != nil {
		return nil, fmt.Errorf("confirm order: %w", err)
	}

	if err := workflow.ExecuteActivity(activityCtx, orderAct.FulfillStock, result.OrderID).Get(ctx, nil); err != nil {
		disconnectedCtx, _ := workflow.NewDisconnectedContext(rollbackCtx)
		cancelCmd := activity.SystemCancelOrderCommand{
			OrderID: result.OrderID,
			Reason:  activity.FulfillStockFailedReason,
		}
		_ = workflow.ExecuteActivity(disconnectedCtx, orderAct.CancelOrderBySystem, cancelCmd).Get(disconnectedCtx, nil)
		_ = workflow.ExecuteActivity(disconnectedCtx, orderAct.ReleaseStock, result.OrderID).Get(disconnectedCtx, nil)
		return nil, fmt.Errorf("fulfill stock: %w", err)
	}

	return &result, nil
}
