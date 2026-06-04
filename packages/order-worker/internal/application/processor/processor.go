package processor

import (
	"context"
	"encoding/json"

	"inventory-worker-module/internal/application/port"
)

type InventoryEventProcessor struct {
	handlers map[string]port.EventHandler
}

func NewInventoryEventProcessor(repo port.ESRepository) port.EventProcessor {
	p := &InventoryEventProcessor{
		handlers: make(map[string]port.EventHandler),
	}

	return p
}

func (p *InventoryEventProcessor) Handle(ctx context.Context, msg port.Message) error {
	h, ok := p.handlers[msg.Topic]
	if !ok {
		return nil
	}

	var rawPayload json.RawMessage = msg.Value
	return h.Handle(ctx, rawPayload)
}
