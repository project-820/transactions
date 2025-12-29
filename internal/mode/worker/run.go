package worker

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/project-820/transactions/internal/adapters/in/eventloop"
	"github.com/project-820/transactions/internal/platform"
	"github.com/project-820/transactions/internal/platform/bootstrap"
	"github.com/project-820/transactions/internal/platform/infra"
	"github.com/project-820/transactions/internal/platform/runner"
)

const httpAddr = ":8080"

func Run(ctx context.Context) error {
	infraMux := infra.NewMux(infra.Params{
		Readiness: nil, // позже: readiness воркера db nats
		Metrics:   nil,
	})

	httpServer, err := bootstrap.NewHTTPServer(&bootstrap.HTTPConfig{Addr: httpAddr}, infraMux)
	if err != nil {
		return fmt.Errorf("failed to init infra http server: %w", err)
	}

	p := platform.NewPlatform(
		runner.NewHTTPRunner(httpAddr, httpServer),
	)

	if err := p.Run(); err != nil {
		return fmt.Errorf("failed to run platform: %w", err)
	}

	eventLoop := eventloop.NewEventLoop()
	go eventLoop.Run(ctx)

	// go syncLoop.Run(ctx)

	<-ctx.Done()
	slog.Info("stop signal received")

	if err := p.Stop(); err != nil {
		return fmt.Errorf("worker stop failed: %w", err)
	}
	return nil
}
