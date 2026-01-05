package repo

import (
	"context"

	"github.com/project-820/transactions/internal/core/usecase"
	"github.com/uptrace/bun"
)

var _ usecase.TxRepo = (*txRepo)(nil)

type txRepo struct {
	db bun.IDB
}

func NewTxRepository(db bun.IDB) *txRepo {
	return &txRepo{
		db: db,
	}
}

func (r *txRepo) InsertOnchainIgnoreConflicts(
	ctx context.Context,
	walletID int64,
	userID string,
	txs []usecase.OnchainTx,
) ([]usecase.InsertedTx, error) {
	return nil, nil
}
