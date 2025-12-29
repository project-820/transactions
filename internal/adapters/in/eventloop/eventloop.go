package eventloop

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/project-820/transactions/internal/core/transactions"
	"github.com/project-820/transactions/pkg/broker"
	"github.com/project-820/transactions/pkg/workerpool"
)

type EventLoop struct {
	workerPool          workerpool.WorkerPool
	walletUpdateUsecase transactions.WalletUpdate
	consumer            broker.Consumer

	log *slog.Logger
}

func NewEventLoop() EventLoop {
	return EventLoop{}
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

		if err := l.workerPool.TaskAdd(
			ctx,
			workerpool.Task{
				F: func(ctx context.Context, data any) error {
					bytes := data.([]byte)
					wallet, err := decodeWalletUpdate(bytes)
					if err != nil {
						// TODO: wrap error
						_ = msg.Nack(ctx, err)

						return fmt.Errorf("decode wallet update: %w", err)
					}

					if err := l.walletUpdateUsecase.Update(ctx, wallet); err != nil {
						// TODO: wrap error
						_ = msg.Nack(ctx, err)

						return fmt.Errorf("update: %w", err)
					}

					// TODO: wrap error
					_ = msg.Ack(ctx)

					return nil
				},
				Data: msg.Payload(),
			},
		); err != nil {
			// TODO: wrap error
			_ = msg.Nack(ctx, err)

			l.log.Error("task add", "reason", err)
		}
	}
}
