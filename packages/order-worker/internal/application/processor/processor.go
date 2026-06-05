package processor

import (
	"context"
	"encoding/json"

	"order-shared-module/alias"
	"order-shared-module/events"
	"order-worker-module/internal/application/port"
	"order-worker-module/internal/application/processor/handler"
)

type OrderEventProcessor struct {
	handlers map[string]port.EventHandler
}

func NewOrderEventProcessor(repo port.ESRepository) port.EventProcessor {
	p := &OrderEventProcessor{
		handlers: make(map[string]port.EventHandler),
	}

	p.handlers[events.TopicOrderCreated] = handler.NewOrderCreatedHandler(repo, alias.OrderAlias)
	p.handlers[events.TopicOrderConfirmed] = handler.NewOrderConfirmedHandler(repo, alias.OrderAlias)
	p.handlers[events.TopicOrderCancelled] = handler.NewOrderCancelledHandler(repo, alias.OrderAlias)

	return p
}

func (p *OrderEventProcessor) Handle(ctx context.Context, msg port.Message) error {
	h, ok := p.handlers[msg.Topic]
	if !ok {
		return nil
	}

	var rawPayload json.RawMessage = msg.Value
	return h.Handle(ctx, rawPayload)
}
