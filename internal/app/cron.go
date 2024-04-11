package app

import (
	"context"

	"github.com/maxuanquang/ojs/internal/handler/jobs"
	"go.uber.org/zap"
)

type Cron struct {
	cron   jobs.Cron
	logger *zap.Logger
}

func NewCron(
	cron jobs.Cron,
	logger *zap.Logger,
) (Cron, error) {
	return Cron{
		cron:   cron,
		logger: logger,
	}, nil
}

func (c *Cron) Start() {
	err := c.cron.Start(context.Background())
	c.logger.With(zap.Error(err)).Error("cron jobs stopped")
}
