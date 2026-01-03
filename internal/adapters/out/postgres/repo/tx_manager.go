package repo

import (
	"context"

	"github.com/project-820/transactions/internal/core/usecase"
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
	fn func(context.Context, usecase.Repositories) error,
) error {
	return m.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		repos := &repositories{tx: tx}
		return fn(ctx, repos)
	})
}

type repositories struct {
	tx bun.Tx
}

func (r *repositories) OutboxWriter() usecase.OutboxWriter {
	return NewOutboxWriter(&r.tx)
}

func (r *repositories) WalletRepo() usecase.WalletRepo {
	return NewWalletRepositiry(&r.tx)
}

func (r *repositories) WalletSyncTaskRepo() usecase.WalletSyncTaskRepo {
	return NewWalletSyncTaskRepository(&r.tx)
}
