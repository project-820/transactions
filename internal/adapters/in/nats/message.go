package nats

import (
	"context"
)

type jsMessage struct {
	data []byte
}

func (m *jsMessage) Data() []byte {
	return m.data
}

func (m *jsMessage) Ack(ctx context.Context) error {
	return nil
}

func (m *jsMessage) Nack(ctx context.Context, reason error) error {
	return nil
}
