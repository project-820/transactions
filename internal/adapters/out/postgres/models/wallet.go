package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Wallet struct {
	bun.BaseModel `bun:"table:wallets,alias:w"`

	ID      int64  `bun:"id,pk,autoincrement"`
	UserID  string `bun:"user_id,notnull"`
	Chain   string `bun:"chain,notnull"`
	Address string `bun:"address,notnull"`
	Label   string `bun:"label,nullzero"`

	// 1=active, 2=disabled
	Status int16 `bun:"status,nullzero,notnull"`

	CreatedAt time.Time `bun:"created_at,notnull,default:now()"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:now()"`
}
