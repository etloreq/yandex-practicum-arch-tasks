package entity

import "time"

type SetStatus struct {
	DeviceID int64
	Enabled  bool
}

type DeviceSettings struct {
	DeviceID  int64
	Enabled   bool
	UpdatedAt time.Time
}
