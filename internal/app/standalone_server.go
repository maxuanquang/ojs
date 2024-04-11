package app

import (
	"context"
	"syscall"

	"github.com/maxuanquang/ojs/internal/handler/consumer"
	"github.com/maxuanquang/ojs/internal/handler/grpc"
	"github.com/maxuanquang/ojs/internal/handler/http"
	"github.com/maxuanquang/ojs/internal/handler/jobs"
	"github.com/maxuanquang/ojs/internal/utils"
	"go.uber.org/zap"
)

type StandaloneServer struct {
	grpcServer grpc.Server
	httpServer http.Server
	mqConsumer consumer.RootConsumer
	cron       jobs.Cron
	logger     *zap.Logger
}

func NewStandaloneServer(
	grpcServer grpc.Server,
	httpServer http.Server,
	mqConsumer consumer.RootConsumer,
	cron jobs.Cron,
	logger *zap.Logger,
) (StandaloneServer, error) {
	return StandaloneServer{
		grpcServer: grpcServer,
		httpServer: httpServer,
		mqConsumer: mqConsumer,
		cron:       cron,
		logger:     logger,
	}, nil
}

func (s *StandaloneServer) Start() {

	go func() {
		err := s.grpcServer.Start(context.Background())
		s.logger.With(zap.Error(err)).Error("gRPC server stopped")
	}()

	go func() {
		err := s.httpServer.Start(context.Background())
		s.logger.With(zap.Error(err)).Error("http server stopped")
	}()

	go func() {
		err := s.mqConsumer.Start(context.Background())
		s.logger.With(zap.Error(err)).Error("mq consumer stopped")
	}()

	go func() {
		err := s.cron.Start(context.Background())
		s.logger.With(zap.Error(err)).Error("cron jobs stopped")
	}()

	utils.WaitForSignals(syscall.SIGINT, syscall.SIGTERM)
}
