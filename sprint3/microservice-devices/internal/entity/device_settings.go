package entity

import "time"

type SetHeating struct {
	DeviceID      int64
	HeatingStatus bool
}

type HeatingSettings struct {
	DeviceID      int64
	HeatingStatus bool
	UpdatedAt     time.Time
}
