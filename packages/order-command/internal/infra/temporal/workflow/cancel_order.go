package workflow

import (
	"fmt"

	"order-command-module/internal/application/commands/cancel_order"
	"order-command-module/internal/application/port"
	"order-command-module/internal/bootstrap/config"
	"order-command-module/internal/infra/temporal/activity"

	"go.temporal.io/sdk/workflow"
)

func CancelOrderWorkflow(ctx workflow.Context, cmd cancel_order.Command, cfg config.TemporalConfig) (*cancel_order.Result, error) {
	activityCtx := activityContext(ctx, cfg)

	var orderAct *activity.OrderActivity
	var result cancel_order.Result
	if err := workflow.ExecuteActivity(activityCtx, orderAct.CancelOrder, cmd).Get(ctx, &result); err != nil {
		return nil, fmt.Errorf("cancel order: %w", err)
	}

	releaseParam := port.ReleaseAndFullfilStockParam{
		OrderID: result.OrderID,
		ActorID: cmd.ActorID.String(),
	}

	if err := workflow.ExecuteActivity(activityCtx, orderAct.ReleaseStock, releaseParam).Get(ctx, nil); err != nil {
		return nil, fmt.Errorf("release stock: %w", err)
	}

	return &result, nil
}
