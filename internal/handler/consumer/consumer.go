package consumer

import (
	"context"
	"encoding/json"

	"github.com/maxuanquang/ojs/internal/dataaccess/mq/consumer"
	"github.com/maxuanquang/ojs/internal/dataaccess/mq/producer"
	"go.uber.org/zap"
)

type RootConsumer interface {
	Start(ctx context.Context) error
}

func NewRootConsumer(
	submissionCreatedHandler SubmissionCreatedHandler,
	mqConsumer consumer.Consumer,
	logger *zap.Logger,
) RootConsumer {
	return &rootConsumer{
		submissionCreatedHandler: submissionCreatedHandler,
		mqConsumer:               mqConsumer,
		logger:                   logger,
	}
}

type rootConsumer struct {
	submissionCreatedHandler SubmissionCreatedHandler
	mqConsumer               consumer.Consumer
	logger                   *zap.Logger
}

// Start implements RootConsumer.
func (r *rootConsumer) Start(ctx context.Context) error {
	r.mqConsumer.RegisterHandler(
		producer.MessageQueueSubmissionCreated,
		func(ctx context.Context, payload []byte) error {

			var submissionID uint64

			err := json.Unmarshal(payload, &submissionID)
			if err != nil {
				return err
			}

			return r.submissionCreatedHandler.Handle(ctx, submissionID)
		},
	)

	return r.mqConsumer.Start(ctx)
}
