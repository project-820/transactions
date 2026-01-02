package nats

import (
	"context"
	"time"

	"github.com/project-820/transactions/pkg/broker"
)

var _ broker.Message = (*message)(nil)

type message struct {
	eventType  string
	payload    []byte
	id         string
	key        string
	occurredAt time.Time

	// raw *nats.Msg или jsMsg, спрятан внутри
}

func (m *message) EventType() string                            { return m.eventType }
func (m *message) Payload() []byte                              { return m.payload }
func (m *message) ID() string                                   { return m.id }
func (m *message) Key() string                                  { return m.key }
func (m *message) OccurredAt() time.Time                        { return m.occurredAt }
func (m *message) Ack(ctx context.Context) error                { return nil }
func (m *message) Nack(ctx context.Context, reason error) error { return nil }
