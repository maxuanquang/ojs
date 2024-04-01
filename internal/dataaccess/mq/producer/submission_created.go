package producer

import (
	"context"
	"encoding/json"

	"github.com/maxuanquang/ojs/internal/utils"
	"go.uber.org/zap"
)

const (
	MessageQueueSubmissionCreated = "submisson_created"
)

type SubmissionCreatedProducer interface {
	Produce(ctx context.Context, submissionID uint64) error
}

func NewSubmissionCreatedProducer(client Client, logger *zap.Logger) (SubmissionCreatedProducer, error) {
	return &submissionCreatedProducer{
		client: client,
		logger: logger,
	}, nil
}

type submissionCreatedProducer struct {
	client Client
	logger *zap.Logger
}

// Produce implements DownloadTaskCreatedProducer.
func (d *submissionCreatedProducer) Produce(ctx context.Context, submissionID uint64) error {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Any("download_task_id", submissionID))

	payload, err := json.Marshal(submissionID)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to marshal event submission created")
		return err
	}

	err = d.client.Produce(ctx, MessageQueueSubmissionCreated, payload)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to produce message submission created")
		return err
	}

	return nil
}
