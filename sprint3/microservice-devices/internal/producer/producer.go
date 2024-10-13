package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/google/uuid"

	"github.com/etloreq/yandex-practicum-arch-tasks/sprint3/microservice-devices/internal/cfg"
	"github.com/etloreq/yandex-practicum-arch-tasks/sprint3/microservice-devices/internal/entity"
)

type producer struct {
	cfg    cfg.Kafka
	client sarama.SyncProducer
}

func New(cfg cfg.Kafka) (*producer, func() error) {
	client, err := sarama.NewSyncProducer([]string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)}, nil)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	return &producer{client: client}, client.Close
}

func (p *producer) SendStatusChangedMessage(_ context.Context, settings entity.SetStatus) error {
	msg := message{
		DeviceID:     settings.DeviceID,
		DeviceStatus: settings.Enabled,
	}

	bytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("build msg: %w", err)
	}

	kafkaMsg := &sarama.ProducerMessage{
		Topic: p.cfg.StatusChangedTopic,
		Key:   sarama.StringEncoder(uuid.New().String()),
		Value: sarama.ByteEncoder(bytes),
	}

	// отправка сообщения в Kafka
	_, _, err = p.client.SendMessage(kafkaMsg)
	if err != nil {
		return fmt.Errorf("send: %w", err)
	}
	return nil
}
