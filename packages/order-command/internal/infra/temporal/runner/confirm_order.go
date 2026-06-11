package runner

import (
	"context"
	"fmt"

	"order-command-module/internal/application/commands/confirm_order"
	"order-command-module/internal/infra/temporal/workflow"

	authx "github.com/iamKienb/go-core/middleware/auth"
	enumspb "go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
)

type ConfirmOrderRunner interface {
	ConfirmOrder(ctx context.Context, cmd confirm_order.Command) (*confirm_order.Result, error)
}

func (r *workflowRunner) ConfirmOrder(ctx context.Context, cmd confirm_order.Command) (*confirm_order.Result, error) {
	requestID := authx.GetRequestID(ctx)
	run, err := r.temporalClient.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
		ID:                       fmt.Sprintf("confirm-order-%s-%s", cmd.OrderID, requestID),
		TaskQueue:                r.temporalCfg.OrderTaskQueue,
		WorkflowIDConflictPolicy: enumspb.WORKFLOW_ID_CONFLICT_POLICY_USE_EXISTING,
		WorkflowIDReusePolicy:    enumspb.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE,
	}, workflow.ConfirmOrderWorkflow, cmd, r.temporalCfg)
	if err != nil {
		return nil, err
	}

	var output confirm_order.Result
	if err := run.Get(ctx, &output); err != nil {
		return nil, fmt.Errorf("confirm order saga: %w", err)
	}
	return &output, nil
}
