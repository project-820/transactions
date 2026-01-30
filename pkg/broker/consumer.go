package broker

import (
	"context"
)

type Message interface {
	Subject() string
	Data() []byte
	Ack() error
	Term() error
	Nak() error
}

type Consumer interface {
	Next(ctx context.Context) (Message, error)
	Close() error
}
