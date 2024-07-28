package kafka

import (
	"github.com/wlcmtunknwndth/messagio_test/backend/internal/config"
	"log/slog"
)

//type Storage interface {
//	Save(ctx context.Context, input *api.Message) (int64, error)
//	ChatMessage(ctx context.Context, id int64, palID int64, limit, offset int) []api.Message
//	Chats(ctx context.Context, id int64) []api.Message
//}

type Kafka struct {
	log     *slog.Logger
	cfg     *config.Broker
	storage Storage
}

const scope = "inner.scope.broker.kafka."

func New(cfg *config.Broker, storage Storage, log *slog.Logger) *Kafka {
	const op = scope + "New"
	//kafkaCfg := kafka.ConfigMap{
	//	"bootstrap.servers": cfg.Servers,
	//	"client.id":         cfg.ClientID,
	//	"acks":              cfg.Acks,
	//	""
	//}

	return &Kafka{
		log:     log,
		cfg:     cfg,
		storage: storage,
	}
}
