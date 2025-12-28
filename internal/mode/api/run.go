package api

import (
	"context"
	"fmt"
	"log/slog"

	platformpkg "github.com/project-820/transactions/internal/platform"
	"github.com/project-820/transactions/internal/platform/bootstrap"
	"github.com/project-820/transactions/internal/platform/infra"
	"github.com/project-820/transactions/internal/platform/runner"
	server "github.com/project-820/transactions/internal/transport/grpc"
)

const (
	grpcAddr = ":8082"
	httpAddr = ":8080"
)

func Run(ctx context.Context) error {
	transactionsServer := server.NewTransactionsServer(server.ServerParams{})

	grpcServer, err := bootstrap.NewGRPCServer(&bootstrap.GRPCConfig{
		Addr:             grpcAddr,
		ValidateMessages: transactionsServer.ValidateMessages,
		ValidatePaths:    transactionsServer.ValidatePaths,
	})
	if err != nil {
		return fmt.Errorf("failed to init grpc server: %w", err)
	}

	infraMux := infra.NewMux(infra.Params{
		Readiness: nil, // позже: DB/NATS readiness
		Metrics:   nil, // nil => expvar
	})

	httpServer, err := bootstrap.NewHTTPServer(&bootstrap.HTTPConfig{Addr: httpAddr}, infraMux)
	if err != nil {
		return fmt.Errorf("failed to init infra http server: %w", err)
	}

	platform := platformpkg.NewPlatform(
		runner.NewHTTPRunner(httpAddr, httpServer),
		runner.NewGRPCRunner(grpcAddr, grpcServer, []runner.GRPCRegisterFunc{
			transactionsServer.Register,
		}),
	)

	if err := platform.Run(); err != nil {
		return fmt.Errorf("failed to run platform: %w", err)
	}

	<-ctx.Done()
	slog.Info("stop signal received")

	if err := platform.Stop(); err != nil {
		return fmt.Errorf("app stop failed: %w", err)
	}
	return nil
}
