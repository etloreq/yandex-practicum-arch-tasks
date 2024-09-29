package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/etloreq/yandex-practicum-arch-tasks/sprint3/microservice-devices/internal/entity"
	"github.com/etloreq/yandex-practicum-arch-tasks/sprint3/microservice-devices/internal/errs"
)

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) BeginTx() (*sqlx.Tx, error) {
	return r.db.Beginx()
}

func (r *repository) EndTx(tx *sqlx.Tx, err error) error {
	if err == nil {
		return tx.Commit()
	}
	return tx.Rollback()
}

func (r *repository) GetHeatingStatus(ctx context.Context, deviceID int64) (entity.HeatingSettings, error) {
	query := `select device_id, heating_status, updated_at from device_settings where device_id = $1`
	var settings entity.HeatingSettings
	err := r.db.QueryRowxContext(ctx, query, deviceID).
		Scan(&settings.DeviceID, &settings.HeatingStatus, &settings.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.HeatingSettings{}, errs.ErrNotFound
	}
	return settings, err
}

func (r *repository) GetHeatingStatusTx(ctx context.Context, tx *sqlx.Tx, deviceID int64) (entity.HeatingSettings, error) {
	query := `select device_id, heating_status, updated_at from device_settings where device_id = $1 for update`
	var settings entity.HeatingSettings
	err := tx.QueryRowxContext(ctx, query, deviceID).
		Scan(&settings.DeviceID, &settings.HeatingStatus, &settings.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.HeatingSettings{}, errs.ErrNotFound
	}
	return settings, err
}

func (r *repository) SetHeatingStatus(ctx context.Context, tx *sqlx.Tx, in entity.SetHeating) error {
	query := `update device_settings set heating_status = $1, updated_at = now() where device_id = $2`
	_, err := tx.ExecContext(ctx, query, in.HeatingStatus, in.DeviceID)
	return err
}

func (r *repository) CreateDefaultSettings(ctx context.Context, tx *sqlx.Tx, deviceID int64) error {
	query := `insert into device_settings(device_id) values($1)`
	_, err := tx.ExecContext(ctx, query, deviceID)
	return err
}
