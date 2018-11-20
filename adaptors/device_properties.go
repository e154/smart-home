package adaptors

import "time"

type DeviceDefault struct {
	Address *int `json:"address"`
}

type DeviceZigBee struct {
	Address string `json:"address"`
}

type DeviceSmartBus struct {
	Address  *int          `json:"address"`
	Baud     int           `json:"baud"`
	Sleep    int64         `json:"sleep"`
	StopBite int64         `json:"stop_bite"`
	Timeout  time.Duration `json:"timeout"`
	Tty      string        `json:"tty"`
}
