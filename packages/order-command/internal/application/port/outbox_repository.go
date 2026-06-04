package port

import (
	"context"
	"time"

	"order-command-module/internal/domain/shared"

	"github.com/google/uuid"
)

type OutboxParam struct {
	AggregateID   uuid.UUID
	AggregateType string
	Events        []shared.DomainEvent
}

type OutboxEvent struct {
	ID             uuid.UUID
	AggregateID    uuid.UUID
	AggregateType  string
	EventType      string
	Payload        interface{}
	PartitionKey   string
	IdempotencyKey uuid.UUID
	CreatedAt      time.Time
}

type OutboxRepository interface {
	SaveOutboxBatch(ctx context.Context, events []OutboxEvent) error
}
