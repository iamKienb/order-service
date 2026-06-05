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

	var orderAct *activity.OrderActivity
	var result confirm_order.Result
	if err := workflow.ExecuteActivity(activityCtx, orderAct.ConfirmOrder, cmd).Get(ctx, &result); err != nil {
		return nil, fmt.Errorf("confirm order: %w", err)
	}

	if err := workflow.ExecuteActivity(activityCtx, orderAct.FulfillStock, result.OrderID).Get(ctx, nil); err != nil {
		return nil, fmt.Errorf("fulfill stock: %w", err)
	}

	return &result, nil
}
