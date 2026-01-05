package usecase

import (
	"context"
	"time"
)

type InsertedTx struct {
	TxID       int64
	WalletID   int64
	UserID     string
	AssetRef   string
	Qty        string
	OccurredAt time.Time
}

type TxRepo interface {
	InsertOnchainIgnoreConflicts(ctx context.Context, walletID int64, userID string, txs []OnchainTx) ([]InsertedTx, error)
}

type WalletSyncTask struct {
	WalletID int64
	UserID   string
	Chain    string
	Address  string

	Cursor string
}

type WalletSyncTaskRepo interface {
	UpsertTasks(ctx context.Context, userID string, walletIDs []int64, runAfter time.Time) error
	ClaimDue(ctx context.Context, limit int, lockTTL time.Duration) ([]WalletSyncTask, error)
	MarkDone(ctx context.Context, walletID int64, nextRunAfter time.Time, nextCursor string) error
	MarkFailed(ctx context.Context, walletID int64, retryAt time.Time, errMsg string) error
}

type OutboxMessage struct {
	Subject   string
	EventName string
	Key       string
	Payload   []byte
}

type OutboxWriterRepo interface {
	AddMessages(ctx context.Context, msgs []OutboxMessage) error
}

type WalletRepo interface {
	UpsertWallets(ctx context.Context, userID string, wallets []WalletRef) ([]int64, error)
}

type Repositories interface {
	TxRepo() TxRepo
	OutboxWriterRepo() OutboxWriterRepo
	WalletRepo() WalletRepo
	WalletSyncTaskRepo() WalletSyncTaskRepo
}

type TxManager interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context, r Repositories) error) error
}
