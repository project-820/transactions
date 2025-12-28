package grpc

import (
	"fmt"
	"log/slog"
	"net"

	"buf.build/go/protovalidate"
	"github.com/project-820/transactions/internal/platform/grpc/interceptors"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type GRPCParams struct {
	ValidateMessages []proto.Message
	ValidatePaths    map[string]bool
	Address          string
}

type GRPC struct {
	*grpc.Server
	listener net.Listener
}

func NewGRPCServer(params *GRPCParams) (*GRPC, error) {
	listener, err := net.Listen("tcp", params.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(),
		protovalidate.WithMessages(params.ValidateMessages...),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize validator: %w", err)
	}

	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			interceptors.MethodScopedValidatioUnaryInterceptor(
				validator,
				params.ValidatePaths,
			),
		),
	)

	return &GRPC{
		Server:   gRPCServer,
		listener: listener,
	}, nil
}

func (s *GRPC) Serve() error {
	slog.Info(fmt.Sprintf("grpc server listening on %s", s.listener.Addr().String()))

	return s.Server.Serve(s.listener)
}

func (s *GRPC) Stop() error {
	slog.Info("shutting down grpc server...")

	s.Server.Stop()

	return nil
}
