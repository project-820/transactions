package ports

import (
	"context"
	"time"
)

type OutboxMessage struct {
	ID          int64
	Subject     string
	EventName   string
	Status      string
	Payload     []byte
	CreatedAt   time.Time
	ProcessedAt *time.Time
}

type OutboxRepository interface {
	AddMessages(ctx context.Context, msgs []OutboxMessage) error
}
