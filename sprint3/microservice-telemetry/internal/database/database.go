package database

import (
	influxdb "github.com/influxdata/influxdb-client-go/v2"

	"github.com/etloreq/yandex-practicum-arch-tasks/sprint3/microservice-telemetry/internal/cfg"
)

func MustDB(cfg cfg.Database) influxdb.Client {
	return influxdb.NewClient(cfg.URL, cfg.Token)
}
