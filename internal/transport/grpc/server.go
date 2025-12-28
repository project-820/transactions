package grpc

import (
	"context"

	pb "github.com/project-820/transactions/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type ServerParams struct {
}

type Server struct {
	ValidateMessages []proto.Message
	ValidatePaths    map[string]bool

	pb.UnimplementedTransactionsServiceServer
}

func NewTransactionsServer(params ServerParams) *Server {
	return &Server{
		ValidateMessages: []proto.Message{
			&pb.HelloRequest{},
		},
		ValidatePaths: map[string]bool{
			"/transactions.AssetsCatalogService/Hello": true,
		},
	}
}

func (s *Server) Register(gs grpc.ServiceRegistrar) {
	pb.RegisterTransactionsServiceServer(gs, s)
}

func (s *Server) Hello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Data: "Response: " + request.GetData(),
	}, nil
}
