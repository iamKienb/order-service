package runner

import (
	"order-command-module/internal/bootstrap/config"

	"go.temporal.io/sdk/client"
)

type Runner interface {
	PlaceOrderRunner
	ConfirmOrderRunner
	CancelOrderRunner
}

type workflowRunner struct {
	temporalClient client.Client
	temporalCfg    config.TemporalConfig
}

func NewWorkflowRunner(temporalClient client.Client, cfg config.TemporalConfig) Runner {
	return &workflowRunner{temporalClient: temporalClient, temporalCfg: cfg}
}
