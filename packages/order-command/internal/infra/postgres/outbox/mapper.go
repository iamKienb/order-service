package outbox

import (
	"encoding/json"

	"inventory-command-module/db/repository"
	"inventory-command-module/internal/application/port"

	"github.com/iamKienb/go-core/postgres/conv"
	"github.com/jackc/pgx/v5/pgtype"
)

func toInfraOutboxBatch(events []port.OutboxEvent) repository.SaveOutboxBatchParams {
	params := repository.SaveOutboxBatchParams{
		Ids:             make([]pgtype.UUID, len(events)),
		AggregateIds:    make([]pgtype.UUID, len(events)),
		AggregateTypes:  make([]string, len(events)),
		EventTypes:      make([]string, len(events)),
		Payloads:        make([][]byte, len(events)),
		PartitionKeys:   make([]string, len(events)),
		IdempotencyKeys: make([]pgtype.UUID, len(events)),
		CreatedAts:      make([]pgtype.Timestamptz, len(events)),
	}

	for i, event := range events {
		payload, err := json.Marshal(event.Payload)
		if err != nil {
			payload = []byte("{}")
		}

		params.Ids[i] = conv.UUID(event.ID)
		params.AggregateIds[i] = conv.UUID(event.AggregateID)
		params.AggregateTypes[i] = event.AggregateType
		params.EventTypes[i] = event.EventType
		params.Payloads[i] = payload
		params.PartitionKeys[i] = event.PartitionKey
		params.IdempotencyKeys[i] = conv.UUID(event.IdempotencyKey)
		params.CreatedAts[i] = conv.TimeStampZ(&event.CreatedAt)
	}

	return params
}
