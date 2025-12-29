package broker

import (
	"context"
	"time"
)

type Consumer interface {
	Next(ctx context.Context) (Message, error)
	Close() error
}

type Message interface {
	EventType() string
	Payload() []byte
	ID() string
	Key() string
	OccurredAt() time.Time

	Ack(ctx context.Context) error
	Nack(ctx context.Context, reason error) error
}
