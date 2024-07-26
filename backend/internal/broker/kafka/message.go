package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/wlcmtunknwndth/messagio_test/common/domain/topics"
)

func (k *Kafka) SaveConsumer() error {
	const op = scope + "SaveConsumer"
	consumer, err := kafka.NewConsumer(k.cfg)
	defer consumer.Close()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	consumer.Subscribe(topics.HandleMessage)
}
