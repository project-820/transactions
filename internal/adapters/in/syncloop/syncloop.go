package syncloop

import (
	"context"
	"log/slog"

	"github.com/project-820/transactions/internal/core/usecase"
	"github.com/project-820/transactions/pkg/workerpool"
)

type SyncLoopParams struct {
	Pool                *workerpool.Pool
	WalletTxSyncUsecase usecase.WalletTxSync

	Log *slog.Logger
}

type SyncLoop struct {
	pool                *workerpool.Pool
	walletTxSyncUsecase usecase.WalletTxSync

	log *slog.Logger
}

func NewSyncLoop(params SyncLoopParams) SyncLoop {
	return SyncLoop{
		pool:                params.Pool,
		walletTxSyncUsecase: params.WalletTxSyncUsecase,
		log:                 params.Log,
	}
}

func (l *SyncLoop) Run(ctx context.Context) {

}

func (l *SyncLoop) Stop() {
	l.pool.StopNow()
}
