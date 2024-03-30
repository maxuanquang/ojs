package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	"github.com/maxuanquang/ojs/internal/configs"
	"github.com/maxuanquang/ojs/internal/generated/grpc/ojs"
	"google.golang.org/grpc"
)

type Server interface {
	Start(ctx context.Context) error
}

func NewServer(grpcConfig configs.GRPC, handler ojs.OjsServiceServer) Server {
	return &server{
		grpcConfig: grpcConfig,
		handler:    handler,
	}
}

type server struct {
	grpcConfig configs.GRPC
	handler    ojs.OjsServiceServer
}

// Start implements Server.
func (s *server) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.grpcConfig.Address)
	if err != nil {
		return err
	}
	defer listener.Close()

	var opts = []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			validator.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			validator.StreamServerInterceptor(),
		),
	}
	server := grpc.NewServer(opts...)
	ojs.RegisterOjsServiceServer(server, s.handler)

	fmt.Printf("gRPC server is running on %s\n", s.grpcConfig.Address)
	return server.Serve(listener)
}
