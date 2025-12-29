package postgres

import (
	"context"
	"fmt"

	"github.com/project-820/transactions/internal/core/ports"

	"github.com/uptrace/bun"
)

var _ ports.OutboxRepository = (*outboxRepository)(nil)

type outboxRepository struct {
	db bun.IDB
}

func NewOutboxRepository(db bun.IDB) *outboxRepository {
	return &outboxRepository{
		db: db,
	}
}

func (r *outboxRepository) AddMessages(ctx context.Context, msgs []ports.OutboxMessage) error {
	_, err := r.db.NewInsert().
		Model(&msgs).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("insert outbox msgs: %w", err)
	}

	return nil
}
