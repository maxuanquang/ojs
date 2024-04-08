package consumer

import (
	"context"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
	"github.com/maxuanquang/ojs/internal/configs"
	"github.com/maxuanquang/ojs/internal/utils"
	"go.uber.org/zap"
)

const (
	consumerClientID = "ojs-consumer"
)

type Consumer interface {
	Start(ctx context.Context) error
	RegisterHandler(queueName string, handlerFunc HandlerFunc)
}

func NewConsumer(
	mqConfig configs.MQ,
	logger *zap.Logger,
) (Consumer, error) {
	saramaConsumerGroup, err := sarama.NewConsumerGroup(mqConfig.Addresses, mqConfig.ConsumerGroupID, newSaramaConfig(mqConfig))
	if err != nil {
		logger.With(zap.Error(err)).Error("can not create sarama consumer group")
		return nil, err
	}

	return &consumer{
		saramaConsumerGroup:       saramaConsumerGroup,
		logger:                    logger,
		queueNameToHandlerFuncMap: make(map[string]HandlerFunc),
	}, nil
}

type HandlerFunc func(ctx context.Context, payload []byte) error

type consumer struct {
	logger                    *zap.Logger
	queueNameToHandlerFuncMap map[string]HandlerFunc
	saramaConsumerGroup       sarama.ConsumerGroup
}

// RegisterHandler implements Consumer.
func (c *consumer) RegisterHandler(queueName string, handlerFunc HandlerFunc) {
	c.queueNameToHandlerFuncMap[queueName] = handlerFunc
}

// Start implements Consumer.
func (c *consumer) Start(ctx context.Context) error {
	logger := utils.LoggerWithContext(ctx, c.logger)

	exitSignalChannel := make(chan os.Signal, 1)
	signal.Notify(exitSignalChannel, os.Interrupt)

	for queueName, handlerFunc := range c.queueNameToHandlerFuncMap {
		go func(queueName string, handlerFunc HandlerFunc) {
			for {
				err := c.saramaConsumerGroup.Consume(ctx, []string{queueName}, newConsumerHandler(handlerFunc, exitSignalChannel))
				if err != nil {
					logger.With(zap.String("queueName", queueName)).With(zap.Error(err)).Error("failed to consume message from queue")
					break
				}
			}
			logger.Info("consumer stopped")
		}(queueName, handlerFunc)
	}

	<-exitSignalChannel
	return nil
}

func newSaramaConfig(_ configs.MQ) *sarama.Config {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.ClientID = consumerClientID
	config.Metadata.Full = true
	return config
}

func newConsumerHandler(
	handlerFunc HandlerFunc,
	exitSignalChannel chan os.Signal,
) sarama.ConsumerGroupHandler {
	return &consumerHandler{
		handlerFunc:       handlerFunc,
		exitSignalChannel: exitSignalChannel,
	}
}

type consumerHandler struct {
	handlerFunc       HandlerFunc
	exitSignalChannel chan os.Signal
}

// ConsumeClaim implements sarama.ConsumerGroupHandler.
func (c *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				session.Commit()
				return nil
			}

			err := c.handlerFunc(session.Context(), message.Value)
			if err != nil {
				return err
			}
		case <-c.exitSignalChannel:
			session.Commit()
			return nil
		}
	}
}

// Cleanup implements sarama.ConsumerGroupHandler.
func (c *consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// Setup implements sarama.ConsumerGroupHandler.
func (c *consumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}
