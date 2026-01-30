package natsjs

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/project-820/transactions/pkg/broker"
)

var _ broker.Consumer = (*jsConsumer)(nil)

type ConsumerOptions struct {
	FilterSubject string
	Name          string // durable name
	AckPolicy     jetstream.AckPolicy
	AckWait       time.Duration
	MaxAckPending int

	FetchMaxWait time.Duration
}

type jsConsumer struct {
	log *slog.Logger

	consumer jetstream.Consumer
	sub      jetstream.MessagesContext
	maxWait  time.Duration
}

func (c *Client) NewConsumer(ctx context.Context, stream jetstream.Stream, opts ConsumerOptions) (broker.Consumer, error) {
	if opts.FetchMaxWait <= 0 {
		opts.FetchMaxWait = 5 * time.Second
	}

	if opts.Name == "" {
		return nil, errors.New("consumer name is required")
	}

	if opts.FilterSubject == "" {
		return nil, errors.New("consumer filter subject is required")
	}

	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		FilterSubject: opts.FilterSubject,
		Name:          opts.Name,
		Durable:       opts.Name,
		AckPolicy:     opts.AckPolicy,
		AckWait:       opts.AckWait,
		MaxAckPending: opts.MaxAckPending,
	})
	if err != nil {
		return nil, fmt.Errorf("create or update consumer: %w", err)
	}

	sub, err := consumer.Messages()
	if err != nil {
		return nil, fmt.Errorf("consumer messages context: %w", err)
	}

	return &jsConsumer{
		consumer: consumer,
		sub:      sub,
		maxWait:  opts.FetchMaxWait,
	}, nil
}

func (c *jsConsumer) Next(ctx context.Context) (broker.Message, error) {
	msgs, err := c.consumer.Fetch(1, jetstream.FetchMaxWait(c.maxWait))
	if err != nil {
		return nil, fmt.Errorf("fetch: %w", err)
	}

	for msg := range msgs.Messages() {
		return jsMessage{msg}, nil
	}

	if err := msgs.Error(); err != nil {
		return nil, fmt.Errorf("fetch result: %w", err)
	}

	return nil, context.DeadlineExceeded
}

func (c *jsConsumer) Close() error {
	if c.sub != nil {
		c.sub.Stop()
	}
	return nil
}

type jsMessage struct {
	msg jetstream.Msg
}

func (m jsMessage) Subject() string { return m.msg.Subject() }
func (m jsMessage) Data() []byte    { return m.msg.Data() }
func (m jsMessage) Ack() error      { return m.msg.Ack() }
func (m jsMessage) Term() error     { return m.msg.Term() }
func (m jsMessage) Nak() error      { return m.msg.Nak() }
