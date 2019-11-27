package endpoint

import "github.com/e154/smart-home/system/mqtt/management"

type MqttEndpoint struct {
	*CommonEndpoint
}

func NewMqttEndpoint(common *CommonEndpoint) *MqttEndpoint {
	return &MqttEndpoint{
		CommonEndpoint: common,
	}
}

func (m *MqttEndpoint) GetClients(limit, offset int) (list []*management.ClientInfo, total int, err error) {
	list, total, err = m.mqtt.Management().GetClients(limit, offset)
	return
}

func (m *MqttEndpoint) GetClient(clientId string) (client *management.ClientInfo, err error) {
	client, err = m.mqtt.Management().GetClient(clientId)
	return
}

func (m *MqttEndpoint) GetSessions(limit, offset int) (list []*management.SessionInfo, total int, err error) {
	list, total, err = m.mqtt.Management().GetSessions(limit, offset)
	return
}

func (m *MqttEndpoint) GetSession(clientId string) (session *management.SessionInfo, err error) {
	session, err = m.mqtt.Management().GetSession(clientId)
	return
}

func (m *MqttEndpoint) GetSubscriptions(clientId string, limit, offset int) (list []*management.SubscriptionInfo, total int, err error) {
	list, total, err = m.mqtt.Management().GetSubscriptions(clientId, limit, offset)
	return
}

func (m *MqttEndpoint) Subscribe(clientId, topic string, qos int) (err error) {
	err = m.mqtt.Management().Subscribe(clientId, topic, qos)
	return
}

func (m *MqttEndpoint) Unsubscribe(clientId, topic string) (err error) {
	err = m.mqtt.Management().Unsubscribe(clientId, topic)
	return
}

func (m *MqttEndpoint) Publish(topic string, qos int, payload []byte, retain bool) (err error) {
	err = m.mqtt.Management().Publish(topic, qos, payload, retain)
	return
}

func (m *MqttEndpoint) CloseClient(clientId string) (err error) {
	err = m.mqtt.Management().CloseClient(clientId)
	return
}
