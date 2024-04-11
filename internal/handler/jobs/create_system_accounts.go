package jobs

import (
	"context"
	"errors"

	"github.com/maxuanquang/ojs/internal/configs"
	"github.com/maxuanquang/ojs/internal/generated/grpc/ojs"
	"github.com/maxuanquang/ojs/internal/logic"
	"go.uber.org/zap"
)

type CreateSystemAccountsJob interface {
	Run(ctx context.Context) error
	GetSchedule() string
}

func NewCreateSystemAccountsJob(
	accountLogic logic.AccountLogic,
	cronConfig configs.Cron,
	logger *zap.Logger,
) (CreateSystemAccountsJob, error) {
	return &createSystemAccountsJob{
		accountLogic: accountLogic,
		cronConfig:   cronConfig,
		logger:       logger,
	}, nil
}

type createSystemAccountsJob struct {
	accountLogic logic.AccountLogic
	cronConfig   configs.Cron
	logger       *zap.Logger
}

// GetSchedule implements CreateSystemAccountsJob.
func (c *createSystemAccountsJob) GetSchedule() string {
	return c.cronConfig.CreateSystemAccounts.Schedule
}

// Run implements CreateSystemAccountsJob.
func (c *createSystemAccountsJob) Run(ctx context.Context) error {
	_, err := c.accountLogic.CreateAccount(
		context.Background(),
		logic.CreateAccountInput{
			Name:     c.cronConfig.CreateSystemAccounts.Admin.Name,
			Password: c.cronConfig.CreateSystemAccounts.Admin.Password,
			Role:     ojs.Role_Admin,
		},
	)
	if err != nil && !errors.Is(err, logic.ErrAccountAlreadyExists) {
		c.logger.Error("create system account failed", zap.Error(err))
		return err
	}
	c.logger.Info("admin account created")

	_, err = c.accountLogic.CreateAccount(
		context.Background(),
		logic.CreateAccountInput{
			Name:     c.cronConfig.CreateSystemAccounts.Worker.Name,
			Password: c.cronConfig.CreateSystemAccounts.Worker.Password,
			Role:     ojs.Role_Worker,
		},
	)
	if err != nil && !errors.Is(err, logic.ErrAccountAlreadyExists) {
		c.logger.Error("create system account failed", zap.Error(err))
		return err
	}
	c.logger.Info("worker account created")

	return nil
}
