package app

import (
	"context"
	"syscall"

	"github.com/maxuanquang/ojs/internal/handler/consumer"
	"github.com/maxuanquang/ojs/internal/utils"
	"go.uber.org/zap"
)

type Worker struct {
	mqConsumer consumer.RootConsumer
	logger     *zap.Logger
}

func NewWorker(
	mqConsumer consumer.RootConsumer,
	logger *zap.Logger,
) (Worker, error) {
	return Worker{
		mqConsumer: mqConsumer,
		logger:     logger,
	}, nil
}

func (w *Worker) Start() {

	go func() {
		err := w.mqConsumer.Start(context.Background())
		w.logger.With(zap.Error(err)).Error("can not start MQ Consumer")
	}()

	utils.WaitForSignals(syscall.SIGINT, syscall.SIGTERM)
}
