package usecase

import (
	"context"
	"time"
)

type OutboxMessage struct {
	Subject   string
	EventName string
	Payload   []byte
}

type OutboxWriter interface {
	AddMessages(ctx context.Context, msgs []OutboxMessage) error
}

type WalletRepo interface {
	UpsertWallets(ctx context.Context, userID string, wallets []WalletRef) ([]int64, error)
}

type Repositories interface {
	OutboxWriter() OutboxWriter
	WalletRepo() WalletRepo
	WalletSyncTaskRepo() WalletSyncTaskRepo
}

type TxManager interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context, r Repositories) error) error
}

type WalletSyncTask struct {
	WalletID int64
	UserID   string
}

type WalletSyncTaskRepo interface {
	UpsertTasks(ctx context.Context, userID string, walletIDs []int64, runAfter time.Time) error
	ClaimDue(ctx context.Context, limit int, lockTTL time.Duration) ([]WalletSyncTask, error)
	MarkDone(ctx context.Context, walletID int64, nextRunAfter time.Time) error
	MarkFailed(ctx context.Context, walletID int64, retryAt time.Time, errMsg string) error
}
