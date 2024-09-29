package producer

import (
	"context"

	"github.com/etloreq/yandex-practicum-arch-tasks/sprint3/microservice-devices/internal/entity"
)

type producer struct {
}

func New() *producer {
	return &producer{}
}

func (p *producer) SendHeatingSettingsChangedMessage(ctx context.Context, settings entity.SetHeating) error {
	// отправка в кафку не реализована
	return nil
}
