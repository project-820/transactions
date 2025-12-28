package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type OutboxMessage struct {
	bun.BaseModel `bun:"table:outbox_messages"`

	ID          int64      `bun:"id,pk,autoincrement"                                    json:"id"`
	Subject     string     `bun:"subject"                                                json:"subject"`
	EventName   string     `bun:"event_name"                                             json:"event_name"`
	Status      string     `bun:"status"                                                 json:"status"`
	Payload     []byte     `bun:"payload"                                                json:"payload"`
	CreatedAt   time.Time  `bun:"created_at,nullzero,notnull,default:current_timestamp"  json:"created_at"`
	ProcessedAt *time.Time `bun:"processed_at,nullzero"                                  json:"processed_at,omitempty"`
}
