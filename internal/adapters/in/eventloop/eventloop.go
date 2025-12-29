package eventloop

import (
	"context"
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

	stopCh chan struct{}
}

func NewEventLoop() EventLoop {
	return EventLoop{
		stopCh: make(chan struct{}),
	}
}

func (l *EventLoop) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			l.log.Error("context", "reason", ctx.Err())
			return
		case <-l.stopCh:
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
				F: func(ctx context.Context, msg any) error {
					wallet := msg.(transactions.Wallet)
					return l.walletUpdateUsecase.Update(ctx, wallet)
				},
				Data: msg,
			},
		); err != nil {
			l.log.Error("task add", "reason", err)
		}
	}
}

func (l *EventLoop) Stop() {
	l.stopCh <- struct{}{}
}
