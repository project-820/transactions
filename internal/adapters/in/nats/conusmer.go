package nats

import (
	"context"
	"log/slog"

	"github.com/project-820/transactions/pkg/broker"
)

type jsConsumer struct {
	log *slog.Logger
}

func NewConsumer() broker.Consumer {
	return &jsConsumer{}
}

func (c *jsConsumer) Next(ctx context.Context) (broker.Message, error) {
	return nil, nil
}

func (c *jsConsumer) Close() error {
	return nil
}
