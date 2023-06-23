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

	Version = "0.0.1"
)

const (
	// AttrSubscribeTopic ...
	AttrSubscribeTopic = "subscribe_topic"
	// AttrMqttLogin ...
	AttrMqttLogin = "mqtt_login"
	// AttrMqttPass ...
	AttrMqttPass = "mqtt_pass"
)

// NewSettings ...
func NewSettings() m.Attributes {
	return m.Attributes{
		AttrSubscribeTopic: {
			Name: AttrSubscribeTopic,
			Type: common.AttributeString,
		},
		AttrMqttLogin: {
			Name: AttrMqttLogin,
			Type: common.AttributeString,
		},
		AttrMqttPass: {
			Name: AttrMqttPass,
			Type: common.AttributeString,
		},
	}
}
