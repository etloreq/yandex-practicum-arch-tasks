package producer

type message struct {
	DeviceID     int64 `json:"device_id"`
	DeviceStatus bool  `json:"device_status"`
}
