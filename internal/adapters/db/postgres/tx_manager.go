package postgres

import (
	"context"

	"github.com/project-820/transactions/internal/services/ports"
	"github.com/uptrace/bun"
)

type txManager struct {
	db *bun.DB
}

func NewTxManager(db *bun.DB) *txManager {
	return &txManager{db: db}
}

func (m *txManager) WithinTx(
	ctx context.Context,
	fn func(context.Context, ports.Repositories) error,
) error {
	return m.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		repos := &repositories{tx: tx}
		return fn(ctx, repos)
	})
}

type repositories struct {
	tx bun.Tx
}

func (r *repositories) Outbox() ports.OutboxRepository {
	return NewOutboxRepository(&r.tx)
}
