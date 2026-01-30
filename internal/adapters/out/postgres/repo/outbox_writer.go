package repo

import (
	"context"
	"fmt"

	"github.com/project-820/transactions/internal/adapters/out/postgres/models"
	"github.com/project-820/transactions/internal/core/usecase"

	"github.com/uptrace/bun"
)

var _ usecase.OutboxWriterRepo = (*outboxWriter)(nil)

type outboxWriter struct {
	db bun.IDB
}

func NewOutboxWriter(db bun.IDB) *outboxWriter {
	return &outboxWriter{
		db: db,
	}
}

func (r *outboxWriter) AddMessages(ctx context.Context, msgs []usecase.OutboxMessage) error {
	modelMsgs := make([]models.OutboxMessage, 0, len(msgs))
	for _, msg := range msgs {
		modelMsgs = append(modelMsgs, models.OutboxMessage{
			Subject:   msg.Subject,
			EventName: msg.EventName,
			Key:       msg.Key,
			Payload:   msg.Payload,
		})
	}

	_, err := r.db.NewInsert().
		Model(&modelMsgs).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("insert outbox msgs: %w", err)
	}

	return nil
}
