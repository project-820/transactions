package worker

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/project-820/transactions/internal/adapters/in/eventloop"
	"github.com/project-820/transactions/internal/adapters/out/natsjs"
	"github.com/project-820/transactions/internal/adapters/out/postgres"
	"github.com/project-820/transactions/internal/adapters/out/postgres/repo"
	"github.com/project-820/transactions/internal/config"
	"github.com/project-820/transactions/internal/core/usecase"
	platformpkg "github.com/project-820/transactions/internal/platform"
	"github.com/project-820/transactions/internal/platform/bootstrap"
	"github.com/project-820/transactions/internal/platform/infra"
	"github.com/project-820/transactions/internal/platform/runner"
	"github.com/project-820/transactions/pkg/workerpool"
)

const httpAddr = ":8080"

func Run(ctx context.Context) error {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		return fmt.Errorf("load config from env: %w", err)
	}
	if err := cfg.ValidateWorker(); err != nil {
		return fmt.Errorf("cfg validate worker: %w", err)
	}

	jsClient, err := natsjs.NewClient(natsjs.Options{
		URL:        cfg.Common.JetStream.Host,
		User:       cfg.Common.JetStream.User,
		Password:   cfg.Common.JetStream.Password,
		ClientName: cfg.Common.JetStream.ClientName,

		Verbose:              cfg.Common.JetStream.Verbose,
		AllowReconnect:       cfg.Common.JetStream.AllowReconnect,
		RetryOnFailedConnect: cfg.Common.JetStream.RetryOnFailedConnect,

		PublishAsyncMaxPending: cfg.Common.JetStream.PublishAsyncMaxPending},
	)
	if err != nil {
		return fmt.Errorf("create natsjs client: %w", err)
	}

	consumer, err := jsClient.NewConsumer(ctx, cfg.Worker.Stream.StreamName,
		natsjs.ConsumerOptions{
			FilterSubject: cfg.Worker.Consumer.Subject,
			Name:          cfg.Worker.Consumer.Name,
			AckPolicy:     cfg.Worker.Consumer.AckPolicy,
			AckWait:       cfg.Worker.Consumer.AckWait,
			MaxAckPending: cfg.Worker.Consumer.MaxAckPending,
		})
	if err != nil {
		return fmt.Errorf("create nats consumer: %w", err)
	}

	db, err := postgres.NewPostgresDB(ctx, postgres.PostgresParams{})
	if err != nil {
		return fmt.Errorf("create postgres db: %w", err)
	}

	txManager := repo.NewTxManager(db)
	pool := workerpool.NewPool(workerpool.Options{
		Workers:   cfg.Worker.Pool.Workers,
		QueueSize: cfg.Worker.Pool.QueueSize,
		OnPanic:   nil,
	})

	eventLoop := eventloop.NewEventLoop(eventloop.EventLoopParams{
		Pool: pool,
		WalletUpdateUsecase: usecase.NewWalletUpdateUsecase(
			usecase.WalletUpdateParams{
				TxManager: txManager,
				Log:       nil,
			},
		),
		Consumer: consumer,
		Log:      nil,
	})
	go eventLoop.Run(ctx)

	// go syncLoop.Run(ctx)

	infraMux := infra.NewMux(infra.Params{
		Readiness: nil, // позже: readiness воркера, db, nats
		Metrics:   nil,
	})

	httpServer, err := bootstrap.NewHTTPServer(bootstrap.HTTPConfig{Addr: httpAddr}, infraMux)
	if err != nil {
		return fmt.Errorf("init infra http server: %w", err)
	}

	platform := platformpkg.NewPlatform(
		runner.NewHTTPRunner(httpAddr, httpServer),
	)

	if err := platform.Run(); err != nil {
		return fmt.Errorf("run platform: %w", err)
	}

	<-ctx.Done()
	slog.Info("stop signal received")
	pool.StopNow()
	_ = consumer.Close()
	jsClient.Close()

	if err := platform.Stop(); err != nil {
		return fmt.Errorf("worker stop: %w", err)
	}

	return nil
}
