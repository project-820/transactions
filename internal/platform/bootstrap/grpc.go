package bootstrap

import (
	"fmt"

	"github.com/project-820/transactions/internal/platform/grpc"
	"google.golang.org/protobuf/proto"
)

type GRPCConfig struct {
	Addr string

	ValidateMessages []proto.Message
	ValidatePaths    map[string]bool
}

func NewGRPCServer(cfg GRPCConfig) (*grpc.GRPC, error) {
	if cfg.Addr == "" {
		return nil, fmt.Errorf("grpc addr is empty")
	}

	grpcServer, err := grpc.NewGRPCServer(
		&grpc.GRPCParams{
			Address:          cfg.Addr,
			ValidateMessages: cfg.ValidateMessages,
			ValidatePaths:    cfg.ValidatePaths,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc server: %w", err)
	}

	return grpcServer, nil
}
