package module

import (
	"inventory-worker-module/internal/application/port"
	"inventory-worker-module/internal/application/processor"
)

type ApplicationModule struct {
	EventProcessor port.EventProcessor
}

func NewApplicationModule(infra *InfraModule) *ApplicationModule {
	return &ApplicationModule{
		EventProcessor: processor.NewInventoryEventProcessor(infra.ESRepo),
	}
}
