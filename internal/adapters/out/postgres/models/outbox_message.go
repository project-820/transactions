package models

import (
	"time"

	"github.com/uptrace/bun"
)

type OutboxMessage struct {
	bun.BaseModel `bun:"table:outbox_messages"`

	ID        int64  `bun:"id,pk,autoincrement"`
	Subject   string `bun:"subject,notnull"`
	EventName string `bun:"event_name,notnull"`
	Key       string `bun:"key,nullzero"`
	Payload   []byte `bun:"payload,notnull,type:bytea"`

	Status      string     `bun:"status,notnull"`
	CreatedAt   time.Time  `bun:"created_at,notnull,default:now()"`
	ProcessedAt *time.Time `bun:"processed_at,nullzero"`

	LockedUntil   *time.Time `bun:"locked_until,nullzero"`
	Attempts      int32      `bun:"attempts,notnull,default:0"`
	NextAttemptAt *time.Time `bun:"next_attempt_at,nullzero"`
	LastError     string     `bun:"last_error,nullzero"`
}
