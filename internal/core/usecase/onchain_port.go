// internal/usecase/onchain_port.go
package usecase

import (
	"context"
	"time"
)

type OnchainClient interface {
	ListTransactions(ctx context.Context, req OnchainTxListRequest) (OnchainTxListResult, error)
}

type OnchainTxListRequest struct {
	Chain   string
	Address string

	Cursor string
	Limit  int
}

type OnchainTxListResult struct {
	Transactions []OnchainTx
	NextCursor   string
	HasMore      bool
}

type OnchainTx struct {
	TxHash     string
	LogIndex   int32
	BlockNum   int64
	OccurredAt time.Time

	Direction string
	AssetRef  string
	Qty       string

	FeeAssetRef string
	FeeQty      string

	Meta []byte
}

type OnchainResolver interface {
	ForChain(chain string) (OnchainClient, bool)
}
