package temporal

import (
	"order-command-module/internal/infra/temporal/activity"
	"order-command-module/internal/infra/temporal/workflow"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

type Worker struct {
	worker worker.Worker
}

type Registry struct {
	OrderActivity *activity.OrderActivity
}

func NewWorker(temporalClient client.Client, taskQueue string, registry Registry) *Worker {
	newWorker := worker.New(temporalClient, taskQueue, worker.Options{
		MaxConcurrentActivityExecutionSize:     20,
		MaxConcurrentWorkflowTaskExecutionSize: 20,
	})

	newWorker.RegisterWorkflow(workflow.PlaceOrderWorkflow)
	newWorker.RegisterWorkflow(workflow.ConfirmOrderWorkflow)
	newWorker.RegisterWorkflow(workflow.CancelOrderWorkflow)
	newWorker.RegisterActivity(registry.OrderActivity)

	return &Worker{worker: newWorker}
}

func (w *Worker) Start() error {
	return w.worker.Start()
}

func (w *Worker) Stop() {
	w.worker.Stop()
}
