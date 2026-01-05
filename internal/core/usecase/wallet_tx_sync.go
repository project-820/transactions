// internal/usecase/wallet_tx_sync.go
package usecase

import (
	"context"
	"fmt"
	"time"
)

type WalletTxSyncParams struct {
	TxManager TxManager
	Onchain   OnchainResolver
	Period    time.Duration
}

type WalletTxSync struct {
	txManager TxManager
	onchain   OnchainResolver
	period    time.Duration
	lockTTL   time.Duration
}

func NewWalletTxSyncUsecase(params WalletTxSyncParams) *WalletTxSync {
	return &WalletTxSync{
		txManager: params.TxManager,
		onchain:   params.Onchain,
		period:    params.Period,
	}
}

func (t *WalletTxSync) Sync(ctx context.Context, task WalletSyncTask) error {
	client, ok := t.onchain.ForChain(task.Chain)
	if !ok {
		return fmt.Errorf("no onchain client for chain=%s", task.Chain)
	}

	onChainTxList, err := client.ListTransactions(ctx, OnchainTxListRequest{
		Chain:   task.Chain,
		Address: task.Address,
		Cursor:  task.Cursor,
		Limit:   500,
	})
	if err != nil {
		retryAt := time.Now().Add(1 * time.Minute)
		_ = t.txManager.WithinTx(ctx, func(ctx context.Context, r Repositories) error {
			return r.WalletSyncTaskRepo().MarkFailed(ctx, task.WalletID, retryAt, err.Error())
		})
		return err
	}

	nextRun := time.Now().Add(t.period)

	return t.txManager.WithinTx(ctx, func(ctx context.Context, r Repositories) error {
		insertedTxs, err := r.TxRepo().InsertOnchainIgnoreConflicts(ctx, task.WalletID, task.UserID, onChainTxList.Transactions)
		if err != nil {
			return err
		}

		msgs := buildTxAddedMessages(insertedTxs)

		if len(msgs) > 0 {
			if err := r.OutboxWriterRepo().AddMessages(ctx, msgs); err != nil {
				return err
			}
		}

		return r.WalletSyncTaskRepo().MarkDone(ctx, task.WalletID, nextRun, onChainTxList.NextCursor)
	})
}

func buildTxAddedMessages(insertedTxs []InsertedTx) []OutboxMessage {
	return nil
}
