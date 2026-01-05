package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/project-820/transactions/internal/adapters/out/postgres/models"
	"github.com/project-820/transactions/internal/core/usecase"
	"github.com/uptrace/bun"
)

var _ usecase.WalletSyncTaskRepo = (*walletSyncTaskRepo)(nil)

type walletSyncTaskRepo struct {
	db bun.IDB
}

func NewWalletSyncTaskRepository(db bun.IDB) *walletSyncTaskRepo {
	return &walletSyncTaskRepo{
		db: db,
	}
}

func (w *walletSyncTaskRepo) UpsertTasks(ctx context.Context, userID string, walletIDs []int64, runAfter time.Time) error {
	if len(walletIDs) == 0 {
		return nil
	}

	walletSyncTasks := make([]models.WalletSyncTask, 0, len(walletIDs))
	for _, walletID := range walletIDs {
		walletSyncTasks = append(walletSyncTasks, models.WalletSyncTask{
			WalletID: walletID,
			UserID:   userID,
			RunAfter: runAfter,
			Status:   1,
		})
	}

	_, err := w.db.NewInsert().
		Model(&walletSyncTasks).
		On("CONFLICT (wallet_id) DO NOTHING").
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("upsert tasks: %w", err)
	}

	return nil
}

func (w *walletSyncTaskRepo) ClaimDue(ctx context.Context, limit int, lockTTL time.Duration) ([]usecase.WalletSyncTask, error) {
	return nil, nil
}

func (w *walletSyncTaskRepo) MarkDone(ctx context.Context, walletID int64, nextRunAfter time.Time, nextCursor string) error {
	return nil
}

func (w *walletSyncTaskRepo) MarkFailed(ctx context.Context, walletID int64, retryAt time.Time, errMsg string) error {
	return nil
}
