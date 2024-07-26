package kafka

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/wlcmtunknwndth/messagio_test/common/domain/api"
	"log/slog"
)

type Storage interface {
	Save(ctx context.Context, input *api.Message) (int64, error)
	ChatMessage(ctx context.Context, id int64, palID int64, limit, offset int) []api.Message
	Chats(ctx context.Context, id int64) []api.Message
}

type Kafka struct {
	log     *slog.Logger
	cfg     *kafka.ConfigMap
	storage Storage
}

const scope = "inner.scope.broker.kafka."

func New(clientId, servers, acks string, storage Storage, log *slog.Logger) *Kafka {
	const op = scope + "New"
	cfg := kafka.ConfigMap{
		"bootstrap.servers": servers,
		"client.id":         clientId,
		"acks":              acks,
	}

	return &Kafka{
		log:     log,
		cfg:     &cfg,
		storage: storage,
	}
}
