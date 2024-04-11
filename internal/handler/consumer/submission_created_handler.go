package consumer

import (
	"context"

	"github.com/maxuanquang/ojs/internal/configs"
	"github.com/maxuanquang/ojs/internal/logic"
	"github.com/maxuanquang/ojs/internal/utils"
	"go.uber.org/zap"
)

type SubmissionCreatedHandler interface {
	Handle(ctx context.Context, submissionID uint64) error
}

func NewSubmissionCreatedHandler(
	accountLogic logic.AccountLogic,
	cronConfig configs.Cron,
	submissionLogic logic.SubmissionLogic,
	logger *zap.Logger,
) (SubmissionCreatedHandler, error) {
	createSessionOutput, err := accountLogic.CreateSession(
		context.Background(),
		logic.CreateSessionInput{
			Name:     cronConfig.CreateSystemAccounts.Worker.Name,
			Password: cronConfig.CreateSystemAccounts.Worker.Password,
		},
	)
	if err != nil {
		return nil, err
	}

	return &submissionCreatedHandler{
		submissionLogic: submissionLogic,
		logger:          logger,
		token:           createSessionOutput.Token,
	}, nil
}

type submissionCreatedHandler struct {
	submissionLogic logic.SubmissionLogic
	logger          *zap.Logger
	token           string
}

// Handle implements DownloadTaskCreatedHandler.
func (d *submissionCreatedHandler) Handle(ctx context.Context, submissionID uint64) error {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Any("submissionID", submissionID))

	logger.Info("submission created event received at handlerFunc")
	err := d.submissionLogic.ExecuteSubmission(
		ctx,
		logic.ExecuteSubmissionInput{
			ID:    submissionID,
			Token: d.token,
		},
	)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to execute submission")
	}

	return nil
}
