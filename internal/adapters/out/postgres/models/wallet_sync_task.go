package models

import (
	"time"

	"github.com/uptrace/bun"
)

type WalletSyncTask struct {
	bun.BaseModel `bun:"table:wallet_sync_tasks"`

	WalletID int64  `bun:"wallet_id,pk"`
	UserID   string `bun:"user_id,notnull"`

	Status int16 `bun:"status,notnull"`

	RunAfter    time.Time  `bun:"run_after,notnull"`
	LockedUntil *time.Time `bun:"locked_until,nullzero"`

	Attempts  int32  `bun:"attempts,notnull"`
	LastError string `bun:"last_error,nullzero"`

	LastStartedAt  *time.Time `bun:"last_started_at,nullzero"`
	LastFinishedAt *time.Time `bun:"last_finished_at,nullzero"`

	CreatedAt time.Time `bun:"created_at,notnull,default:now()"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:now()"`
}
