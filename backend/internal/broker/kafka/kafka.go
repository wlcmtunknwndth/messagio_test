package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	KafkaLib "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/wlcmtunknwndth/messagio_test/backend/internal/config"
	"github.com/wlcmtunknwndth/messagio_test/common/domain/api"
	"github.com/wlcmtunknwndth/messagio_test/common/domain/topics"
	"log/slog"
)

type Kafka struct {
	log            *slog.Logger
	messageCounter *KafkaLib.Producer
}

const scope = "inner.internal.broker.kafka."

func New(cfg *config.Broker, log *slog.Logger) (*Kafka, error) {
	const op = scope + "New"

	counter, err := KafkaLib.NewProducer(&KafkaLib.ConfigMap{
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
	err = k.messageCounter.Produce(&KafkaLib.Message{
		TopicPartition: KafkaLib.TopicPartition{Topic: &topic, Partition: KafkaLib.PartitionAny},
		Value:          data,
	}, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (k *Kafka) Close() error {
	k.messageCounter.Close()
	return nil
}
