package entity

import "time"

type Telemetry struct {
	DeviceID  int64
	Measure   int64
	Timestamp time.Time
}
