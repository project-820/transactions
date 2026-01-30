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
	l.pool.Start(ctx)

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

		subject := msg.Subject()

		task := func(ctx context.Context) error {
			bytes := msg.Data()
			wallet, err := decodeWalletUpdate(bytes)
			if err != nil {
				l.log.Error("decode wallet update failed",
					"subject", subject,
					"err", err,
				)

				if err := msg.Term(); err != nil {
					l.log.Error("msg term failed", "subject", subject, "err", err)
				}

				return fmt.Errorf("decode wallet update: %w", err)
			}

			if err := l.walletUpdateUsecase.Update(ctx, wallet); err != nil {
				l.log.Error("wallet update usecase failed",
					"subject", subject,
					"user_id", wallet.UserID,
					"wallets_n", len(wallet.Wallets),
					"err", err,
				)

				if err := msg.Nak(); err != nil {
					l.log.Error("msg nak failed", "subject", subject, "err", err)
				}

				return fmt.Errorf("wallet update: %w", err)
			}

			// if ackErr := msg.Ack(); ackErr != nil {
			// 	l.log.Error("msg ack failed",
			// 		"subject", subject,
			// 		"user_id", wallet.UserID,
			// 		"err", ackErr,
			// 	)

			// 	return fmt.Errorf("ack: %w", ackErr)
			// }

			l.log.Info("wallet update processed",
				"subject", subject,
				"user_id", wallet.UserID,
				"wallets_n", len(wallet.Wallets),
			)

			return nil
		}

		if err := l.pool.Submit(ctx, task); err != nil {
			l.log.Error("pool submit failed", "subject", subject, "err", err)

			if nakErr := msg.Nak(); nakErr != nil {
				l.log.Error("msg nak failed", "subject", subject, "err", nakErr)
			}
		}
	}
}

func (l *EventLoop) Stop() {
	l.pool.StopNow()
}
