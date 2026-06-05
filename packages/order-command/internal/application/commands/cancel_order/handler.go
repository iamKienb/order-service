package cancel_order

import "context"

type workflowRunner interface {
	CancelOrder(ctx context.Context, cmd Command) (*Result, error)
}

type handler struct {
	workflow workflowRunner
}

func NewHandler(workflow workflowRunner) Executor {
	return &handler{workflow: workflow}
}

func (h *handler) Execute(ctx context.Context, cmd Command) (*Result, error) {
	return h.workflow.CancelOrder(ctx, cmd)
}
