package ports

import "context"

type TxManager interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context, r Repositories) error) error
}

type Repositories interface {
	Outbox() OutboxRepository
}
