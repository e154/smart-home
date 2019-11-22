package mqtt

// Javascript Binding
//
// mqtt
//	.publish
//
type MqttBind struct {
	mqtt *Mqtt
}

func NewMqttBind(mqtt *Mqtt) *MqttBind {
	return &MqttBind{mqtt: mqtt}
}

func (m MqttBind) Publish(topic string, payload []byte, qos uint8, retain bool) {
	if m.mqtt.server == nil {
		return
	}
	m.mqtt.server.Publish(topic, payload, qos, retain)
}
