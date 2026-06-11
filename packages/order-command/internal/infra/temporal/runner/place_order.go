package runner

import (
	"context"
	"fmt"

	"order-command-module/internal/application/commands/place_order"
	"order-command-module/internal/infra/temporal/workflow"

	authx "github.com/iamKienb/go-core/middleware/auth"
	enumspb "go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
)

type PlaceOrderRunner interface {
	PlaceOrder(ctx context.Context, cmd place_order.Command) (*place_order.Result, error)
}

func (r *workflowRunner) PlaceOrder(ctx context.Context, cmd place_order.Command) (*place_order.Result, error) {
	workflowID := placeOrderWorkflowID(ctx, cmd)
	run, err := r.temporalClient.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
		ID:                       workflowID,
		TaskQueue:                r.temporalCfg.OrderTaskQueue,
		WorkflowIDConflictPolicy: enumspb.WORKFLOW_ID_CONFLICT_POLICY_USE_EXISTING,
		WorkflowIDReusePolicy:    enumspb.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE,
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

func placeOrderWorkflowID(ctx context.Context, cmd place_order.Command) string {
	requestID := authx.GetRequestID(ctx)
	return fmt.Sprintf("place-order-%s-%s-%s", cmd.BuyerID.String(), cmd.ShopID.String(), requestID)
}
