package workflow

import (
	"fmt"

	"order-command-module/internal/application/commands/confirm_order"
	"order-command-module/internal/application/port"
	"order-command-module/internal/bootstrap/config"
	"order-command-module/internal/infra/temporal/activity"

	"go.temporal.io/sdk/workflow"
)

func ConfirmOrderWorkflow(ctx workflow.Context, cmd confirm_order.Command, cfg config.TemporalConfig) (*confirm_order.Result, error) {
	activityCtx := activityContext(ctx, cfg)

	var orderAct *activity.OrderActivity
	var result confirm_order.Result

	fullfilParam := port.ReleaseAndFullfilStockParam{
		OrderID: cmd.OrderID,
		ActorID: cmd.ActorID.String(),
	}

	fmt.Println("fullfilParam", fullfilParam)

	if err := workflow.ExecuteActivity(activityCtx, orderAct.FulfillStock, fullfilParam).Get(ctx, nil); err != nil {
		return nil, fmt.Errorf("fulfill stock: %w", err)
	}

	if err := workflow.ExecuteActivity(activityCtx, orderAct.ConfirmOrder, cmd).Get(ctx, &result); err != nil {
		return nil, fmt.Errorf("confirm order: %w", err)
	}

	return &result, nil
}
