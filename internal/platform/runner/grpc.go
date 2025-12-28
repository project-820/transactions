package runner

import (
	"fmt"
	"log"
	"log/slog"

	grpcserver "github.com/project-820/transactions/internal/platform/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var _ Runner = (*GRPCRunner)(nil)

type GRPCRegisterFunc func(gs grpc.ServiceRegistrar)

type GRPCRunner struct {
	grpcAddr       string
	grpcServer     *grpcserver.GRPC
	grpcRegistrars []GRPCRegisterFunc
}

func NewGRPCRunner(
	grpcAddr string,
	grpcServer *grpcserver.GRPC,
	grpcRegistrars []GRPCRegisterFunc,
) *GRPCRunner {
	return &GRPCRunner{
		grpcAddr:       grpcAddr,
		grpcServer:     grpcServer,
		grpcRegistrars: grpcRegistrars,
	}
}

func (g *GRPCRunner) Start() error {
	g.registerGRPCServices()

	if err := g.startGRPCServer(); err != nil {
		return fmt.Errorf("failed to start grpc server: %w", err)
	}

	return nil
}

func (g *GRPCRunner) Stop() error {
	if err := g.grpcServer.Stop(); err != nil {
		return fmt.Errorf("failed to stop gRPC server: %w", err)
	}

	slog.Info("gRPC server stopped")

	return nil
}

func (g *GRPCRunner) startGRPCServer() error {
	go func() {
		if err := g.grpcServer.Serve(); err != nil {
			log.Fatalf("failed to serve grpc server: %v", err)
		}
	}()

	slog.Info(fmt.Sprintf("serving grpc-server on http://localhost%s", g.grpcAddr))

	return nil
}

func (g *GRPCRunner) registerGRPCServices() {
	for _, grpcRegistrar := range g.grpcRegistrars {
		grpcRegistrar(g.grpcServer)
	}

	reflection.Register(g.grpcServer)
}
