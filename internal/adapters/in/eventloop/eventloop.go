package eventloop

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/project-820/transactions/internal/core/usecase"
	"github.com/project-820/transactions/pkg/broker"
	"github.com/project-820/transactions/pkg/workerpool"
)

type EventLoopParams struct {
	Pool                *workerpool.Pool
	WalletUpdateUsecase usecase.WalletUpdate
	Consumer            broker.Consumer

	Log *slog.Logger
}

type EventLoop struct {
	pool                *workerpool.Pool
	walletUpdateUsecase usecase.WalletUpdate
	consumer            broker.Consumer

	log *slog.Logger
}

func NewEventLoop(params EventLoopParams) EventLoop {
	return EventLoop{
		pool:                params.Pool,
		walletUpdateUsecase: params.WalletUpdateUsecase,
		consumer:            params.Consumer,
		log:                 params.Log,
	}
}

func (l *EventLoop) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			l.log.Error("context", "reason", ctx.Err())
			return
		default:
		}

		msg, err := l.consumer.Next(ctx)
		if err != nil {
			l.log.Error("next message", "reason", err)
			continue
		}

		task := func(ctx context.Context) error {
			bytes := msg.Data()
			wallet, err := decodeWalletUpdate(bytes)
			if err != nil {
				// TODO: wrap error, log
				_ = msg.Nak()

				return fmt.Errorf("decode wallet update: %w", err)
			}

			if err := l.walletUpdateUsecase.Update(ctx, wallet); err != nil {
				// TODO: wrap error, log
				_ = msg.Nak()

				return fmt.Errorf("update: %w", err)
			}

			// TODO: wrap error, log
			_ = msg.Ack()

			return nil
		}

		if err := l.pool.Submit(ctx, task); err != nil {
			// TODO: wrap error, log
			_ = msg.Nak()

			l.log.Error("task add", "reason", err)
		}
	}
}
