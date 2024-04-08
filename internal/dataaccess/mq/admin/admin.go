package admin

import (
	"context"
	"errors"

	"github.com/IBM/sarama"
	"github.com/maxuanquang/ojs/internal/configs"
	"go.uber.org/zap"
)

type Admin interface {
	Setup(ctx context.Context) error
}

func NewAdmin(
	logger *zap.Logger,
	mqConfig configs.MQ,
) (Admin, error) {
	clusterAdmin, err := sarama.NewClusterAdmin(
		mqConfig.Addresses,
		sarama.NewConfig(),
	)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create cluster admin")
		return nil, err
	}

	return &admin{
		logger:       logger,
		mqConfig:     mqConfig,
		clusterAdmin: clusterAdmin,
	}, nil
}

type admin struct {
	clusterAdmin sarama.ClusterAdmin
	mqConfig     configs.MQ
	logger       *zap.Logger
}

// Setup implements Broker.
func (b *admin) Setup(ctx context.Context) error {

	err := b.clusterAdmin.CreateTopic(
		b.mqConfig.Topic,
		&sarama.TopicDetail{
			NumPartitions:     int32(b.mqConfig.NumPartitions),
			ReplicationFactor: 1,
		},
		false)
	if err != nil {
		if errors.Is(err, sarama.ErrTopicAlreadyExists) {
			b.logger.Info("topic already exists")
		} else {
			b.logger.With(zap.Error(err)).Error("failed to create topic")
			return err
		}
	}

	err = b.clusterAdmin.CreatePartitions(
		b.mqConfig.Topic,
		int32(b.mqConfig.NumPartitions),
		make([][]int32, 0),
		false,
	)
	if err != nil {
		b.logger.With(zap.Error(err)).Warn("failed to set number of partitions")
	}

	return nil
}
