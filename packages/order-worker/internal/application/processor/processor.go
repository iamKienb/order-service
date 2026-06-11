package processor

import (
	"context"
	"fmt"
	"time"

	"order-shared-module/alias"
	"order-shared-module/events"
	"order-worker-module/internal/application/port"
	"order-worker-module/internal/application/processor/handler"
)

const (
	idemKeyTTL = 24 * time.Hour
	key        = "user-worker:key:%s"
)

type OrderEventProcessor struct {
	handlers    map[string]port.EventHandler
	workerCache port.WorkerCache
}

func NewOrderEventProcessor(repo port.ESRepository, workerCache port.WorkerCache) port.EventProcessor {
	p := &OrderEventProcessor{
		handlers:    make(map[string]port.EventHandler),
		workerCache: workerCache,
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

	idemKey := msg.IdempotencyKey()

	if idemKey != "" {
		key := fmt.Sprintf(key, idemKey)
		isNew, err := p.workerCache.SetNx(ctx, key, 1, idemKeyTTL)
		if err != nil {
			return err
		}

		if !isNew {
			return nil
		}
	}

	return h.Handle(ctx, msg.Value)
}
