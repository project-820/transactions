package repo

import (
	"context"
	"fmt"

	"github.com/project-820/transactions/internal/core/usecase"

	"github.com/uptrace/bun"
)

var _ usecase.OutboxWriter = (*outboxWriter)(nil)

type outboxWriter struct {
	db bun.IDB
}

func NewOutboxWriter(db bun.IDB) *outboxWriter {
	return &outboxWriter{
		db: db,
	}
}

func (r *outboxWriter) AddMessages(ctx context.Context, msgs []usecase.OutboxMessage) error {
	_, err := r.db.NewInsert().
		Model(&msgs).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("insert outbox msgs: %w", err)
	}

	return nil
}
