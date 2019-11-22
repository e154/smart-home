package devices

import (
	. "github.com/e154/smart-home/common"
)

const (
	DevTypeMqtt = DeviceType("mqtt")
)

type DevMqttConfig struct {
	Validation
	User     string `json:"user"`
	Password string `json:"password"`
}

type DevMqttRequest struct {
	Topic   string `json:"topic"`
	Payload []byte `json:"payload"`
	Qos     uint8  `json:"qos"`
	Retain  bool   `json:"retain"`
}
