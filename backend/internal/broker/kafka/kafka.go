package kafka

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log/slog"
)

type Storage interface {
	Send(ctx context.Context)
}

type Kafka struct {
	log *slog.Logger
	cfg *kafka.ConfigMap
}

const scope = "inner.scope.broker.kafka."

func New(clientId, servers, acks string, log *slog.Logger) *Kafka {
	const op = scope + "New"

	//producer, err := kafka.NewProducer()
	//
	//if err != nil {
	//	return nil, fmt.Errorf("%s: %w", op, err)
	//}
	cfg := kafka.ConfigMap{
		"bootstrap.servers": servers,
		"client.id":         clientId,
		"acks":              acks,
	}
	return &Kafka{
		log: log,
		cfg: &cfg,
	}
}
