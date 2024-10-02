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

func (r *repository) GetStatus(ctx context.Context, deviceID int64) (entity.DeviceSettings, error) {
	query := `select device_id, enabled, updated_at from device_settings where device_id = $1`
	var settings entity.DeviceSettings
	err := r.db.QueryRowxContext(ctx, query, deviceID).
		Scan(&settings.DeviceID, &settings.Enabled, &settings.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.DeviceSettings{}, errs.ErrNotFound
	}
	return settings, err
}

func (r *repository) GetStatusTx(ctx context.Context, tx *sqlx.Tx, deviceID int64) (entity.DeviceSettings, error) {
	query := `select device_id, enabled, updated_at from device_settings where device_id = $1 for update`
	var settings entity.DeviceSettings
	err := tx.QueryRowxContext(ctx, query, deviceID).
		Scan(&settings.DeviceID, &settings.Enabled, &settings.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.DeviceSettings{}, errs.ErrNotFound
	}
	return settings, err
}

func (r *repository) SetStatus(ctx context.Context, tx *sqlx.Tx, in entity.SetStatus) error {
	query := `update device_settings set enabled = $1, updated_at = now() where device_id = $2`
	_, err := tx.ExecContext(ctx, query, in.Enabled, in.DeviceID)
	return err
}

func (r *repository) CreateDefaultSettings(ctx context.Context, tx *sqlx.Tx, deviceID int64) error {
	query := `insert into device_settings(device_id) values($1)`
	_, err := tx.ExecContext(ctx, query, deviceID)
	return err
}
