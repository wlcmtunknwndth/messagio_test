package kafka

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/wlcmtunknwndth/messagio_test/backend/internal/config"
	"log/slog"
)

type Kafka struct {
	log *slog.Logger
	//cfg *config.Broker
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

func (k *Kafka) CountMessageSent(ctx context.Context, id, palID, createdAtUnix int64) error {
	const op = scope + "CountMessageSent"

	//k.messageCounter.Produce(&kafka.Message{
	//	TopicPartition: ,
	//
	//})

	return nil
}
