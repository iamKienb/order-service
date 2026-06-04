package outbox

import (
	"context"
	"fmt"

	"inventory-command-module/internal/application/port"
)

func (r *repositoryImpl) SaveOutboxBatch(ctx context.Context, events []port.OutboxEvent) error {
	if err := r.getQuerier(ctx).SaveOutboxBatch(ctx, toInfraOutboxBatch(events)); err != nil {
		return fmt.Errorf("infra: save outbox batch: %w", err)
	}

	return nil
}
