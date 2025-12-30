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

type WalletUpdate struct {
	walletsRepo WalletRepo

	log *slog.Logger
}

func NewWalletUpdate() WalletUpdate {
	return WalletUpdate{}
}

func (w *WalletUpdate) Update(ctx context.Context, wallet WalletUpdateEvent) error {
	_, err := w.walletsRepo.UpsertWallets(ctx, wallet.UserID, wallet.Wallets)
	if err != nil {
		return fmt.Errorf("upsert wallets: %w", err)
	}

	return nil
}
