// internal/usecase/onchain_port.go
package usecase

import (
	"context"
	"time"
)

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

type OnchainClient interface {
	ListTransactions(ctx context.Context, req OnchainTxListRequest) (OnchainTxListResult, error)
}

type OnchainResolver interface {
	ForChain(chain string) (OnchainClient, bool)
}
