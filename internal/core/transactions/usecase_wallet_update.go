package transactions

import (
	"context"
	"log/slog"

	"github.com/project-820/transactions/internal/core/ports"
)

type WalletUpdate struct {
	txManager ports.TxManager

	log *slog.Logger
}

func NewWalletUpdate() WalletUpdate {
	return WalletUpdate{}
}

func (w *WalletUpdate) Update(ctx context.Context, wallet WalletUpdateEvent) error {
	return nil
}
