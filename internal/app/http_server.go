package app

import (
	"context"
	"syscall"

	"github.com/maxuanquang/ojs/internal/handler/grpc"
	"github.com/maxuanquang/ojs/internal/handler/http"
	"github.com/maxuanquang/ojs/internal/utils"
	"go.uber.org/zap"
)

type HTTPServer struct {
	grpcServer grpc.Server
	httpServer http.Server
	logger     *zap.Logger
}

func NewHTTPServer(
	grpcServer grpc.Server,
	httpServer http.Server,
	logger *zap.Logger,
) (HTTPServer, error) {
	return HTTPServer{
		grpcServer: grpcServer,
		httpServer: httpServer,
		logger:     logger,
	}, nil
}

func (h *HTTPServer) Start() {

	go func() {
		err := h.grpcServer.Start(context.Background())
		h.logger.With(zap.Error(err)).Error("gRPC server stopped")
	}()

	go func() {
		err := h.httpServer.Start(context.Background())
		h.logger.With(zap.Error(err)).Error("http server stopped")
	}()

	utils.WaitForSignals(syscall.SIGINT, syscall.SIGTERM)
}
