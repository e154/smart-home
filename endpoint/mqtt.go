// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package endpoint

import (
	"errors"
	"fmt"
	"github.com/e154/smart-home/common/debug"
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

func (m *MqttEndpoint) GetClients(limit, offset int) (list []interface{}, total int, err error) {
	if m.mqtt.Admin() == nil {
		err = ErrMqttServerNoWorked
		return
	}
	q, w, err := m.mqtt.Admin().GetClients(uint32(limit), uint32(offset))
	fmt.Println("-------")
	debug.Println(q)
	debug.Println(w)
	fmt.Println(err)
	return
}

func (m *MqttEndpoint) GetClient(clientId string) (client interface{}, err error) {
	//if m.mqtt.Management() == nil {
	//	err = ErrMqttServerNoWorked
	//	return
	//}
	//client, err = m.mqtt.Management().GetClient(clientId)
	return
}

func (m *MqttEndpoint) GetSessions(limit, offset int) (list []interface{}, total int, err error) {
	//if m.mqtt.Management() == nil {
	//	err = ErrMqttServerNoWorked
	//	return
	//}
	//list, total, err = m.mqtt.Management().GetSessions(limit, offset)
	return
}

func (m *MqttEndpoint) GetSession(clientId string) (session interface{}, err error) {
	//if m.mqtt.Management() == nil {
	//	err = ErrMqttServerNoWorked
	//	return
	//}
	//session, err = m.mqtt.Management().GetSession(clientId)
	return
}

func (m *MqttEndpoint) GetSubscriptions(clientId string, limit, offset int) (list []interface{}, total int, err error) {
	//if m.mqtt.Management() == nil {
	//	err = ErrMqttServerNoWorked
	//	return
	//}
	//list, total, err = m.mqtt.Management().GetSubscriptions(clientId, limit, offset)
	return
}

func (m *MqttEndpoint) Subscribe(clientId, topic string, qos int) (err error) {
	//if m.mqtt.Management() == nil {
	//	err = ErrMqttServerNoWorked
	//	return
	//}
	//err = m.mqtt.Management().Subscribe(clientId, topic, qos)
	return
}

func (m *MqttEndpoint) Unsubscribe(clientId, topic string) (err error) {
	//if m.mqtt.Management() == nil {
	//	err = ErrMqttServerNoWorked
	//	return
	//}
	//err = m.mqtt.Management().Unsubscribe(clientId, topic)
	return
}

func (m *MqttEndpoint) Publish(topic string, qos int, payload []byte, retain bool) (err error) {
	//if m.mqtt.Management() == nil {
	//	err = ErrMqttServerNoWorked
	//	return
	//}
	//err = m.mqtt.Management().Publish(topic, qos, payload, retain)
	return
}

func (m *MqttEndpoint) CloseClient(clientId string) (err error) {
	//if m.mqtt.Management() == nil {
	//	err = ErrMqttServerNoWorked
	//	return
	//}
	//err = m.mqtt.Management().CloseClient(clientId)
	return
}

func (m *MqttEndpoint) SearchTopic(query string, limit, offset int) (result []interface{}, total int64, err error) {
	//if m.mqtt.Management() == nil {
	//	err = ErrMqttServerNoWorked
	//	return
	//}
	//
	//if result, err = m.mqtt.Management().SearchTopic(query); err != nil {
	//	return
	//}
	//
	//// add custom text as topic
	//var exist bool
	//for _, sub := range result {
	//	if sub.Name == query {
	//		exist = true
	//	}
	//}
	//
	//if !exist {
	//	sub := &management.SubscriptionInfo{
	//		Name: query,
	//	}
	//	result = append(result, sub)
	//}

	return
}
