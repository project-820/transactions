package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type WalletUpdateEvent struct {
	EventID    string
	OccurredAt time.Time
	UserID     string
	Wallets    []WalletRef
}

type WalletRef struct {
	Chain   string
	Address string
	Label   string
	Status  int16
}

type WalletUpdateParams struct {
	TxManager TxManager
	Log       *slog.Logger
}

type WalletUpdate struct {
	txManager TxManager
	log       *slog.Logger
}

func NewWalletUpdateUsecase(params WalletUpdateParams) WalletUpdate {
	return WalletUpdate{
		txManager: params.TxManager,
		log:       params.Log,
	}
}

func (w *WalletUpdate) Update(ctx context.Context, wallet WalletUpdateEvent) error {
	err := w.txManager.WithinTx(ctx, func(ctx context.Context, r Repositories) error {
		walletIDs, err := r.WalletRepo().UpsertWallets(ctx, wallet.UserID, wallet.Wallets)
		if err != nil {
			return fmt.Errorf("upsert wallets: %w", err)
		}

		if err := r.WalletSyncTaskRepo().UpsertTasks(ctx, wallet.UserID, walletIDs, time.Now()); err != nil {
			return fmt.Errorf("upsert tasks: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("tx manager within tx: %w", err)
	}

	return nil
}
