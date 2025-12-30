package usecase

import (
	"context"
)

type OutboxMessage struct {
	Subject   string
	EventName string
	Payload   []byte
}

type OutboxWriter interface {
	AddMessages(ctx context.Context, msgs []OutboxMessage) error
}

type WalletRepo interface {
	UpsertWallets(ctx context.Context, userID string, wallets []WalletRef) ([]int64, error)
}

type Repositories interface {
	OutboxWriter() OutboxWriter
	Wallet() WalletRepo
}

type TxManager interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context, r Repositories) error) error
}
