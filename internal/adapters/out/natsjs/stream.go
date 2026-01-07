package natsjs

import (
	"context"
	"errors"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
)

func (c *Client) NewStream(ctx context.Context, cfg jetstream.StreamConfig) (jetstream.Stream, error) {
	stream, err := c.JS().Stream(ctx, cfg.Name)
	if err != nil {
		if errors.Is(err, jetstream.ErrStreamNotFound) {
			if stream, err = c.JS().CreateStream(ctx, cfg); err != nil {
				return nil, fmt.Errorf("create stream: %w:", err)
			}
		} else {
			return nil, fmt.Errorf("get stream %q: %w", cfg.Name, err)
		}
	}

	return stream, nil
}
