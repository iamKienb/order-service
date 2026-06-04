package port

import (
	"context"
	"encoding/json"

	kafkax "github.com/iamKienb/go-core/kafka"
)

type Message = kafkax.Message

type EventProcessor interface {
	Handle(ctx context.Context, msg Message) error
}

type EventHandler interface {
	Handle(ctx context.Context, payload json.RawMessage) error
}
