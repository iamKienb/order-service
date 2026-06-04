package outbox

import (
	"context"

	"order-command-module/internal/application/port"
)

type Service interface {
	PublishBatch(ctx context.Context, params []port.OutboxParam) error
}

type outboxService struct {
	outboxRepo port.OutboxRepository
}

func NewOutboxService(outboxRepo port.OutboxRepository) Service {
	return &outboxService{outboxRepo: outboxRepo}
}
