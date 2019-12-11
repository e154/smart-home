package mqtt_client

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type Subscribe struct {
	Qos      byte
	Callback MQTT.MessageHandler
}
