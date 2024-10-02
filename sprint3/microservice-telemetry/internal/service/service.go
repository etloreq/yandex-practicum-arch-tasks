package service

import (
	"context"

	"github.com/etloreq/yandex-practicum-arch-tasks/sprint3/microservice-telemetry/internal/entity"
)

type repository interface {
	SaveTelemetry(ctx context.Context, telemetry entity.Telemetry) error
	GetLatest(ctx context.Context, deviceID int64) (entity.Telemetry, error)
}

type service struct {
	repo repository
}

func New(repo repository) *service {
	return &service{repo: repo}
}

func (s *service) SaveTelemetry(ctx context.Context, telemetry entity.Telemetry) error {
	return s.repo.SaveTelemetry(ctx, telemetry)
}

func (s *service) GetLatest(ctx context.Context, deviceID int64) (entity.Telemetry, error) {
	return s.repo.GetLatest(ctx, deviceID)
}
