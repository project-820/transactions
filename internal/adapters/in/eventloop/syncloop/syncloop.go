package syncloop

import (
	"log/slog"

	"github.com/project-820/transactions/pkg/workerpool"
)

type SyncLoopParams struct {
	Pool *workerpool.Pool

	Log *slog.Logger
}

type SyncLoop struct {
	Pool *workerpool.Pool

	Log *slog.Logger
}

func NewSyncLoop(params SyncLoopParams) SyncLoop {
	return SyncLoop{
		Pool: params.Pool,
		Log:  params.Log,
	}
}
