package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/wlcmtunknwndth/messagio_test/common/domain/api"
	"github.com/wlcmtunknwndth/messagio_test/common/sl"
	"github.com/wlcmtunknwndth/stats/internal/config"
	"log/slog"
)

type Storage interface {
	CountMessage(ctx context.Context, message *api.Message) error
}

type Kafka struct {
	log            *slog.Logger
	messageHandler *kafka.Consumer
	storage        Storage
}

const scope = "stats.internal.kafka."

func New(cfg *config.Broker, storage Storage, log *slog.Logger) (*Kafka, error) {
	const op = scope + "New"

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Servers,
		"client.id":         cfg.ClientID,
		"acks":              cfg.Acks,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Kafka{
		log:            log,
		messageHandler: consumer,
		storage:        storage,
	}, nil
}

func (k *Kafka) FetchMessages(ctx context.Context) {
	const op = scope + "FetchMessages"

	for {
		select {
		case <-ctx.Done():
			k.log.Error("Context cancelled", sl.Op(op))
			return
		default:
			ev := k.messageHandler.Poll(250)
			switch e := ev.(type) {
			case *kafka.Message:
				_, err := k.messageHandler.CommitMessage(e)
				if err != nil {
					k.log.Error("Error commiting message", sl.Op(op), sl.Err(err))
					continue
				}

				var msg api.Message
				if err = json.Unmarshal(e.Value, &msg); err != nil {
					k.log.Error("Couldn't unmarshal message", sl.Op(op), sl.Err(err))
					continue
				}

				if err = k.storage.CountMessage(ctx, &msg); err != nil {
					k.log.Error("Couldn't count message", sl.Op(op), sl.Err(err))
					continue
				}
			case kafka.PartitionEOF:
				k.log.Warn("Reached limit", sl.Op(op))
				continue
			case kafka.Error:
				k.log.Error("Kafka error", sl.Op(op), sl.Err(e))
				return
				//default:
				//	k.log.Info("Ignored")
			}
		}
	}
}

func (k *Kafka) Close() error {
	const op = scope + "Close"
	if err := k.messageHandler.Close(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
