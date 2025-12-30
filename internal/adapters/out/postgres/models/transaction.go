package models

import (
	"encoding/json"
	"time"

	"github.com/uptrace/bun"
)

type Transaction struct {
	bun.BaseModel `bun:"table:transactions,alias:t"`

	ID int64 `bun:"id,pk,autoincrement"`

	WalletID int64  `bun:"wallet_id,notnull"`
	UserID   string `bun:"user_id,notnull"`

	// source: "onchain" | "cex" | "manual"
	Source string `bun:"source,notnull"`

	// Для onchain:
	TxHash   string `bun:"tx_hash,nullzero"`
	LogIndex int32  `bun:"log_index,nullzero"` // 0..N, для уникальности внутри tx_hash

	// Для cex/manual:
	ExternalID string `bun:"external_id,nullzero"` // id операции у провайдера/ручной

	AssetRef string `bun:"asset_ref,notnull"` // например "eth:0x...", "btc", "usdt:tron", ...
	Qty      string `bun:"qty,notnull"`       // numeric в БД, но в Go часто удобнее string/decimal (см. ниже)

	PriceUSD string `bun:"price_usd,nullzero"`
	FeeAsset string `bun:"fee_asset_ref,nullzero"`
	FeeQty   string `bun:"fee_qty,nullzero"`

	OccurredAt time.Time `bun:"occurred_at,notnull"`

	Meta json.RawMessage `bun:"meta,type:jsonb,nullzero"`

	CreatedAt time.Time `bun:"created_at,notnull,default:now()"`

	// Wallet *Wallet `bun:"rel:belongs-to,join:wallet_id=id"`
}
