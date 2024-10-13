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
	GetStatus(ctx context.Context, deviceID int64) (entity.DeviceSettings, error)
	GetStatusTx(ctx context.Context, tx *sqlx.Tx, deviceID int64) (entity.DeviceSettings, error)
	SetStatus(ctx context.Context, tx *sqlx.Tx, in entity.SetStatus) error
	CreateDefaultSettings(ctx context.Context, tx *sqlx.Tx, deviceID int64) error
}

type producer interface {
	SendStatusChangedMessage(ctx context.Context, settings entity.SetStatus) error
}

type service struct {
	repo     repository
	producer producer
}

func NewService(repo repository, producer producer) *service {
	return &service{repo: repo, producer: producer}
}

func (s *service) GetStatus(ctx context.Context, deviceID int64) (entity.DeviceSettings, error) {
	return s.repo.GetStatus(ctx, deviceID)
}

func (s *service) SetStatus(ctx context.Context, in entity.SetStatus) error {
	tx, err := s.repo.BeginTx()
	if err != nil {
		return fmt.Errorf("start tx: %w", err)
	}

	_, err = s.repo.GetStatusTx(ctx, tx, in.DeviceID)
	if err != nil {
		_ = s.repo.EndTx(tx, err)
		return err
	}

	err = s.repo.SetStatus(ctx, tx, in)
	if err != nil {
		_ = s.repo.EndTx(tx, err)
		return fmt.Errorf("save to db: %w", err)
	}

	if err = s.repo.EndTx(tx, nil); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	err = s.producer.SendStatusChangedMessage(ctx, in)
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

	_, err = s.repo.GetStatusTx(ctx, tx, deviceID)
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
