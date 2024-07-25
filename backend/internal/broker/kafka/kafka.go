package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log/slog"
)

type Kafka struct {
	log      *slog.Logger
	producer *kafka.Producer
}

const scope = "inner.scope.broker.kafka."

func New(clientId, servers, acks string, log *slog.Logger) (*Kafka, error) {
	const op = scope + "New"

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": servers,
		"client.id":         clientId,
		"acks":              acks,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Kafka{
		log:      log,
		producer: producer,
	}, nil
}

func (k *Kafka) Close() {
	k.producer.Close()
}
