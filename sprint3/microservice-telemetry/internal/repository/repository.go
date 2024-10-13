package repository

import (
	"context"
	"fmt"
	"strconv"

	influxdb "github.com/influxdata/influxdb-client-go/v2"

	"github.com/etloreq/yandex-practicum-arch-tasks/sprint3/microservice-telemetry/internal/entity"
	"github.com/etloreq/yandex-practicum-arch-tasks/sprint3/microservice-telemetry/internal/errs"
)

type repository struct {
	client influxdb.Client
	org    string
	bucket string
}

func New(client influxdb.Client, org, bucket string) *repository {
	return &repository{
		client: client,
		org:    org,
		bucket: bucket,
	}
}

func (r *repository) SaveTelemetry(ctx context.Context, telemetry entity.Telemetry) error {
	writeAPI := r.client.WriteAPIBlocking(r.org, r.bucket)

	p := influxdb.NewPointWithMeasurement("telemetry").
		AddTag("device_id", strconv.FormatInt(telemetry.DeviceID, 10)).
		AddField("measure", telemetry.Measure).
		SetTime(telemetry.Timestamp)

	err := writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return writeAPI.Flush(ctx)
}

func (r *repository) GetLatest(ctx context.Context, deviceID int64) (entity.Telemetry, error) {
	query := fmt.Sprintf(`from(bucket: "telemetry")
  |> range(start: -1h)
  |> filter(fn: (r) => r["_measurement"] == "telemetry")
  |> filter(fn: (r) => r["device_id"] == "%s")
  |> last()`, strconv.FormatInt(deviceID, 10))

	queryAPI := r.client.QueryAPI(r.org)
	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		return entity.Telemetry{}, err
	}

	telemetry := entity.Telemetry{
		DeviceID: deviceID,
	}

	if result.Err() != nil {
		return entity.Telemetry{}, fmt.Errorf("result: %w", err)
	}

	if result.Next() {
		value, ok := result.Record().Value().(int64)
		if !ok {
			return entity.Telemetry{}, fmt.Errorf("invalid value: %v", result.Record().Value())
		}
		telemetry.Timestamp = result.Record().Time()
		telemetry.Measure = value
		return telemetry, nil
	}

	return entity.Telemetry{}, errs.ErrNotFound
}
