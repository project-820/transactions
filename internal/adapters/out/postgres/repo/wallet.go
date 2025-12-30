package postgres

import (
	"context"
	"fmt"

	"github.com/project-820/transactions/internal/adapters/out/postgres/models"
	"github.com/project-820/transactions/internal/core/usecase"
	"github.com/uptrace/bun"
)

var _ usecase.WalletRepo = (*wallet)(nil)

type wallet struct {
	db bun.IDB
}

func NewWalletRepositiry(db bun.IDB) *wallet {
	return &wallet{
		db: db,
	}
}

func (r *wallet) UpsertWallets(ctx context.Context, userID string, wallets []usecase.WalletRef) (newWalletIDs []int64, err error) {
	walletModels := make([]models.Wallet, 0, len(wallets))
	for _, wallet := range wallets {
		walletModels = append(walletModels, models.Wallet{
			UserID:  userID,
			Chain:   wallet.Chain,
			Address: wallet.Address,
			Label:   wallet.Label,
			Status:  wallet.Status,
		})
	}

	var rows []models.Wallet
	err = r.db.NewInsert().
		Model(&walletModels).
		On("CONFLICT (user_id, chain, address) DO UPDATE").
		Set("label = EXCLUDED.label").
		Set("status = EXCLUDED.status").
		Set("updated_at = now()").
		Returning("id, user_id, chain, address, status").
		Scan(ctx, &rows)

	if err != nil {
		return nil, fmt.Errorf("insert wallets: %w", err)
	}

	return nil, nil
}
