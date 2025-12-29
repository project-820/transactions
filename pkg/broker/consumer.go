package broker

import "context"

type Consumer interface {
	Next(ctx context.Context) (Message, error)
	Close() error
}

type Message interface {
	Data() []byte

	Ack(ctx context.Context) error

	Nack(ctx context.Context, reason error) error
}
