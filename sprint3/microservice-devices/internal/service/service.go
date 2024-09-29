package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/etloreq/yandex-practicum-arch-tasks/sprint3/microservice-devices/internal/entity"
	"github.com/etloreq/yandex-practicum-arch-tasks/sprint3/microservice-devices/internal/errs"
)

type repository interface {
	BeginTx() (*sqlx.Tx, error)
	EndTx(tx *sqlx.Tx, err error) error
	GetHeatingStatus(ctx context.Context, deviceID int64) (entity.HeatingSettings, error)
	GetHeatingStatusTx(ctx context.Context, tx *sqlx.Tx, deviceID int64) (entity.HeatingSettings, error)
	SetHeatingStatus(ctx context.Context, tx *sqlx.Tx, in entity.SetHeating) error
	CreateDefaultSettings(ctx context.Context, tx *sqlx.Tx, deviceID int64) error
}

type producer interface {
	SendHeatingSettingsChangedMessage(ctx context.Context, settings entity.SetHeating) error
}

type service struct {
	repo     repository
	producer producer
}

func NewService(repo repository, producer producer) *service {
	return &service{repo: repo, producer: producer}
}

func (s *service) GetHeatingStatus(ctx context.Context, deviceID int64) (entity.HeatingSettings, error) {
	return s.repo.GetHeatingStatus(ctx, deviceID)
}

func (s *service) SetHeatingStatus(ctx context.Context, in entity.SetHeating) error {
	tx, err := s.repo.BeginTx()
	if err != nil {
		return fmt.Errorf("start tx: %w", err)
	}

	_, err = s.repo.GetHeatingStatusTx(ctx, tx, in.DeviceID)
	if err != nil {
		_ = s.repo.EndTx(tx, err)
		return err
	}

	err = s.repo.SetHeatingStatus(ctx, tx, in)
	if err != nil {
		_ = s.repo.EndTx(tx, err)
		return fmt.Errorf("save to db: %w", err)
	}

	if err = s.repo.EndTx(tx, nil); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	err = s.producer.SendHeatingSettingsChangedMessage(ctx, in)
	if err != nil {
		return fmt.Errorf("send msg: %w", err)
	}
	return nil
}

func (s *service) RegisterDevice(ctx context.Context, deviceID int64) error {
	tx, err := s.repo.BeginTx()
	if err != nil {
		return fmt.Errorf("start tx: %w", err)
	}

	_, err = s.repo.GetHeatingStatusTx(ctx, tx, deviceID)
	if err != errs.ErrNotFound {
		_ = s.repo.EndTx(tx, err)
		return errors.New("already exists")
	}

	err = s.repo.CreateDefaultSettings(ctx, tx, deviceID)
	if err != nil {
		_ = s.repo.EndTx(tx, err)
		return fmt.Errorf("save to db: %w", err)
	}

	if err = s.repo.EndTx(tx, nil); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}
