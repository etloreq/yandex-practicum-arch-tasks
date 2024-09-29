package entity

import "time"

type Telemetry struct {
	DeviceID    int64
	Temperature int64
	Timestamp   time.Time
}
