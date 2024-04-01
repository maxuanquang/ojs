package consumer

import (
	"context"

	"github.com/maxuanquang/ojs/internal/logic"
	"github.com/maxuanquang/ojs/internal/utils"
	"go.uber.org/zap"
)

type SubmissionCreatedHandler interface {
	Handle(ctx context.Context, submissionID uint64) error
}

func NewSubmissionCreatedHandler(
	submissionLogic logic.SubmissionLogic,
	logger *zap.Logger,
) (SubmissionCreatedHandler, error) {
	return &submissionCreatedHandler{
		submissionLogic: submissionLogic,
		logger:          logger,
	}, nil
}

type submissionCreatedHandler struct {
	submissionLogic logic.SubmissionLogic
	logger          *zap.Logger
}

// Handle implements DownloadTaskCreatedHandler.
func (d *submissionCreatedHandler) Handle(ctx context.Context, submissionID uint64) error {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Any("submissionID", submissionID))

	logger.Info("submission created event received at handlerFunc")
	// err := d.submissionLogic.ExecuteSubmission(
	// 	ctx,
	// 	logic.ExecuteSubmissionInput{},
	// )

	// if err != nil {
	// 	logger.With(zap.Error(err)).Error("failed to download event")
	// 	return err
	// }

	return nil
}
