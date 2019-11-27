package mqtt

import "github.com/e154/smart-home/system/mqtt/management"

type IManagement interface {
	GetClients(limit, offset int) (list []*management.ClientInfo, total int, err error)
	GetClient(clientId string) (client *management.ClientInfo, err error)
	GetSessions(limit, offset int) (list []*management.SessionInfo, total int, err error)
	GetSession(clientId string) (session *management.SessionInfo, err error)
	GetSubscriptions(clientId string, limit, offset int) (list []*management.SubscriptionInfo, total int, err error)
	Subscribe(clientId, topic string, qos int) (err error)
	Unsubscribe(clientId, topic string) (err error)
	Publish(topic string, qos int, payload []byte, retain bool) (err error)
	CloseClient(clientId string) (err error)
}
