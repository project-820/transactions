package repo

import (
	"context"
	"fmt"

	"github.com/project-820/transactions/internal/adapters/out/postgres/models"
	"github.com/project-820/transactions/internal/core/usecase"
	"github.com/uptrace/bun"
)

var _ usecase.WalletRepo = (*walletRepo)(nil)

type walletRepo struct {
	db bun.IDB
}

func NewWalletRepository(db bun.IDB) *walletRepo {
	return &walletRepo{
		db: db,
	}
}

func (r *walletRepo) UpsertWallets(ctx context.Context, userID string, wallets []usecase.WalletRef) ([]int64, error) {
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

	var rows []int64
	err := r.db.NewInsert().
		Model(&walletModels).
		On("CONFLICT (user_id, chain, address) DO UPDATE").
		Set("label = EXCLUDED.label").
		Set("status = EXCLUDED.status").
		Set("updated_at = now()").
		Returning("id").
		Scan(ctx, &rows)

	if err != nil {
		return nil, fmt.Errorf("upsert wallets: %w", err)
	}

	return rows, nil
}
