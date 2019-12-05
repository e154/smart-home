package endpoint

import (
	"errors"
	"github.com/e154/smart-home/system/mqtt/management"
)

var (
	ErrMqttServerNoWorked = errors.New("mqtt server not worked")
)

type MqttEndpoint struct {
	*CommonEndpoint
}

func NewMqttEndpoint(common *CommonEndpoint) *MqttEndpoint {
	return &MqttEndpoint{
		CommonEndpoint: common,
	}
}

func (m *MqttEndpoint) GetClients(limit, offset int) (list []*management.ClientInfo, total int, err error) {
	if m.mqtt.Management() == nil {
		err = ErrMqttServerNoWorked
		return
	}
	list, total, err = m.mqtt.Management().GetClients(limit, offset)
	return
}

func (m *MqttEndpoint) GetClient(clientId string) (client *management.ClientInfo, err error) {
	if m.mqtt.Management() == nil {
		err = ErrMqttServerNoWorked
		return
	}
	client, err = m.mqtt.Management().GetClient(clientId)
	return
}

func (m *MqttEndpoint) GetSessions(limit, offset int) (list []*management.SessionInfo, total int, err error) {
	if m.mqtt.Management() == nil {
		err = ErrMqttServerNoWorked
		return
	}
	list, total, err = m.mqtt.Management().GetSessions(limit, offset)
	return
}

func (m *MqttEndpoint) GetSession(clientId string) (session *management.SessionInfo, err error) {
	if m.mqtt.Management() == nil {
		err = ErrMqttServerNoWorked
		return
	}
	session, err = m.mqtt.Management().GetSession(clientId)
	return
}

func (m *MqttEndpoint) GetSubscriptions(clientId string, limit, offset int) (list []*management.SubscriptionInfo, total int, err error) {
	if m.mqtt.Management() == nil {
		err = ErrMqttServerNoWorked
		return
	}
	list, total, err = m.mqtt.Management().GetSubscriptions(clientId, limit, offset)
	return
}

func (m *MqttEndpoint) Subscribe(clientId, topic string, qos int) (err error) {
	if m.mqtt.Management() == nil {
		err = ErrMqttServerNoWorked
		return
	}
	err = m.mqtt.Management().Subscribe(clientId, topic, qos)
	return
}

func (m *MqttEndpoint) Unsubscribe(clientId, topic string) (err error) {
	if m.mqtt.Management() == nil {
		err = ErrMqttServerNoWorked
		return
	}
	err = m.mqtt.Management().Unsubscribe(clientId, topic)
	return
}

func (m *MqttEndpoint) Publish(topic string, qos int, payload []byte, retain bool) (err error) {
	if m.mqtt.Management() == nil {
		err = ErrMqttServerNoWorked
		return
	}
	err = m.mqtt.Management().Publish(topic, qos, payload, retain)
	return
}

func (m *MqttEndpoint) CloseClient(clientId string) (err error) {
	if m.mqtt.Management() == nil {
		err = ErrMqttServerNoWorked
		return
	}
	err = m.mqtt.Management().CloseClient(clientId)
	return
}

func (m *MqttEndpoint) SearchTopic(query string, limit, offset int) (result []*management.SubscriptionInfo, total int64, err error) {
	if m.mqtt.Management() == nil {
		err = ErrMqttServerNoWorked
		return
	}
	result, err = m.mqtt.Management().SearchTopic(query)
	return
}
