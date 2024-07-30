package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/wlcmtunknwndth/messagio_test/backend/internal/config"
	"github.com/wlcmtunknwndth/messagio_test/common/domain/api"
	"github.com/wlcmtunknwndth/messagio_test/common/domain/topics"
	"log/slog"
)

type Kafka struct {
	log            *slog.Logger
	messageCounter *kafka.Producer
}

const scope = "inner.scope.broker.kafka."

func New(cfg *config.Broker, log *slog.Logger) (*Kafka, error) {
	const op = scope + "New"

	counter, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Servers,
		"client.id":         cfg.ClientID,
		"acks":              cfg.Acks,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Kafka{
		log:            log,
		messageCounter: counter,
	}, nil
}

func (k *Kafka) CountMessageSent(ctx context.Context, msg *api.Message) error {
	const op = scope + "CountMessageSent"

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	topic := topics.HandleCountMessage
	err = k.messageCounter.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          data,
	}, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
