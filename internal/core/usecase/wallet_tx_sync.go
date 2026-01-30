// internal/usecase/wallet_tx_sync.go
package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

const (
	SubjectTxAddedBatch = "tx.added.batch"
	EventTxAddedBatch   = "TxAddedBatch"

	txLimit = 500
)

type WalletTxSyncParams struct {
	TxManager  TxManager
	Onchain    OnchainResolver
	SyncPeriod time.Duration
}

type WalletTxSync struct {
	txManager       TxManager
	onchainResolver OnchainResolver
	syncPeriod      time.Duration
	lockTTL         time.Duration
}

func NewWalletTxSyncUsecase(params WalletTxSyncParams) WalletTxSync {
	return WalletTxSync{
		txManager:       params.TxManager,
		onchainResolver: params.Onchain,
		syncPeriod:      params.SyncPeriod,
	}
}

func (t *WalletTxSync) Sync(ctx context.Context, task WalletSyncTask) error {
	client, ok := t.onchainResolver.ForChain(task.Chain)
	if !ok {
		return fmt.Errorf("no onchain client for chain=%s", task.Chain)
	}

	onChainTxList, err := client.ListTransactions(ctx, OnchainTxListRequest{
		Chain:   task.Chain,
		Address: task.Address,
		Cursor:  task.Cursor,
		Limit:   txLimit,
	})
	if err != nil {
		retryAt := time.Now().Add(1 * time.Minute)
		_ = t.txManager.WithinTx(ctx, func(ctx context.Context, r Repositories) error {
			return r.WalletSyncTaskRepo().MarkFailed(ctx, task.WalletID, retryAt, err.Error())
		})
		return err
	}

	nextRun := time.Now().Add(t.syncPeriod)

	return t.txManager.WithinTx(ctx, func(ctx context.Context, r Repositories) error {
		insertedTxs, err := r.TxRepo().InsertOnchainIgnoreConflicts(ctx, task.WalletID, task.UserID, onChainTxList.Transactions)
		if err != nil {
			return err
		}

		msgs, err := buildTxAddedMessages(insertedTxs)
		if err != nil {
			return err
		}
		if len(msgs) > 0 {
			if err := r.OutboxWriterRepo().AddMessages(ctx, msgs); err != nil {
				return err
			}
		}

		return r.WalletSyncTaskRepo().MarkDone(ctx, task.WalletID, nextRun, onChainTxList.NextCursor)
	})
}

func buildTxAddedMessages(insertedTxs []InsertedTx) ([]OutboxMessage, error) {
	if len(insertedTxs) == 0 {
		return nil, nil
	}

	txsByWalletID := make(map[int64][]InsertedTx, len(insertedTxs))
	for _, insertedTx := range insertedTxs {
		batch := txsByWalletID[insertedTx.WalletID]
		batch = append(batch, insertedTx)
		txsByWalletID[insertedTx.WalletID] = batch
	}

	msgs := make([]OutboxMessage, 0, len(insertedTxs))
	for _, insertedTxsBatch := range txsByWalletID {
		userID := insertedTxsBatch[0].UserID

		payload, err := json.Marshal(insertedTxsBatch)
		if err != nil {
			return nil, fmt.Errorf("marshall inserted tx batch: %w", err)
		}

		msgs = append(msgs, OutboxMessage{
			Subject:   SubjectTxAddedBatch,
			EventName: EventTxAddedBatch,
			Key:       userID,
			Payload:   payload,
		})
	}

	return msgs, nil
}
