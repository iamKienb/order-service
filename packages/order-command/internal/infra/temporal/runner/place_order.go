package runner

import (
	"context"
	"fmt"

	"order-command-module/internal/application/commands/place_order"
	"order-command-module/internal/infra/temporal/workflow"

	"go.temporal.io/sdk/client"
)

type PlaceOrderRunner interface {
	PlaceOrder(ctx context.Context, cmd place_order.Command) (*place_order.Result, error)
}

func (r *workflowRunner) PlaceOrder(ctx context.Context, cmd place_order.Command) (*place_order.Result, error) {
	run, err := r.temporalClient.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
		ID:        fmt.Sprintf("place-order-%s-%s", cmd.BuyerID.String(), cmd.ShopID.String()),
		TaskQueue: r.temporalCfg.OrderTaskQueue,
	}, workflow.PlaceOrderWorkflow, cmd, r.temporalCfg)
	if err != nil {
		return nil, err
	}

	var output place_order.Result
	if err := run.Get(ctx, &output); err != nil {
		return nil, fmt.Errorf("place order saga: %w", err)
	}
	return &output, nil
}
