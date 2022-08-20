package mqtt

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

const (
	// Name ...
	Name = "mqtt"
	// EntityMqtt ...
	EntityMqtt = string("mqtt")
)

const (
	// FuncMqttEvent ...
	FuncMqttEvent = "mqttEvent"
	// FuncEntityAction ...
	FuncEntityAction = "entityAction"
)

const (
	// AttrSubscribeTopic ...
	AttrSubscribeTopic = "subscribe_topic"
)

// NewSettings ...
func NewSettings() m.Attributes {
	return m.Attributes{
		AttrSubscribeTopic: {
			Name: AttrSubscribeTopic,
			Type: common.AttributeString,
		},
	}
}
