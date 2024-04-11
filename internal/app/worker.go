package app

import (
	"context"

	"github.com/maxuanquang/ojs/internal/handler/consumer"
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

	err := w.mqConsumer.Start(context.Background())
	w.logger.With(zap.Error(err)).Error("mq consumer stopped")

}
