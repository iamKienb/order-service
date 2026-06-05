package module

import (
	"order-worker-module/internal/application/port"
	"order-worker-module/internal/application/processor"
)

type ApplicationModule struct {
	EventProcessor port.EventProcessor
}

func NewApplicationModule(infra *InfraModule) *ApplicationModule {
	return &ApplicationModule{
		EventProcessor: processor.NewOrderEventProcessor(infra.ESRepo),
	}
}
