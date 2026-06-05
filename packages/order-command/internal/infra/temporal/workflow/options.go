package workflow

import (
	"time"

	"order-command-module/internal/bootstrap/config"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func activityContext(ctx workflow.Context, cfg config.TemporalConfig) workflow.Context {
	return workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:              cfg.OrderTaskQueue,
		ScheduleToCloseTimeout: cfg.ActivityTimeout,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2,
			MaximumInterval:    10 * time.Second,
			MaximumAttempts:    3,
		},
	})
}

func rollbackContext(ctx workflow.Context, cfg config.TemporalConfig) workflow.Context {
	return workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:              cfg.OrderTaskQueue,
		ScheduleToCloseTimeout: cfg.RollbackTimeout,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2,
			MaximumInterval:    10 * time.Second,
			MaximumAttempts:    5,
		},
	})
}
