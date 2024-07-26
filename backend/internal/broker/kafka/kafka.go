package kafka

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/wlcmtunknwndth/messagio_test/common/domain/api"
	"log/slog"
)

type Storage interface {
	Send(ctx context.Context, input *api.Message)
}

type Kafka struct {
	log *slog.Logger
	cfg *kafka.ConfigMap
}

const scope = "inner.scope.broker.kafka."

func New(clientId, servers, acks string, log *slog.Logger) *Kafka {
	const op = scope + "New"
	cfg := kafka.ConfigMap{
		"bootstrap.servers": servers,
		"client.id":         clientId,
		"acks":              acks,
	}

	//producer, err := kafka.NewProducer(&cfg)
	//if err != nil {
	//	//return nil, fmt.Errorf("%s: %w", op, err)
	//}
	////producer.
	//consumer, err := kafka.NewConsumer(&cfg)
	//consumer.
	return &Kafka{
		log: log,
		cfg: &cfg,
	}
}
