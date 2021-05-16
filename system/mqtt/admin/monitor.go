// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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

package admin

import (
	"container/list"
	"github.com/DrmagicE/gmqtt"
	"github.com/DrmagicE/gmqtt/config"
	"github.com/DrmagicE/gmqtt/persistence/subscription"
	"github.com/DrmagicE/gmqtt/server"
	"github.com/e154/smart-home/common"
	"strings"
	"sync"
	"time"
)

const (
	// Online ...
	Online = "online"
	// Offline ...
	Offline = "offline"
)

type monitor struct {
	clientMu       sync.Mutex
	clientList     *quickList
	subMu          sync.Mutex
	subscriptions  map[string]*quickList // key by clientID
	config         config.Config
	subStatsReader subscription.StatsReader
	statsReader    server.StatsReader
}

// newMonitor
func newMonitor(subStatsReader subscription.StatsReader,
	config config.Config,
	statsReader server.StatsReader) *monitor {
	return &monitor{
		clientList:     newQuickList(),
		subscriptions:  make(map[string]*quickList),
		subStatsReader: subStatsReader,
		config:         config,
		statsReader:    statsReader,
	}
}

// addClient
func (m *monitor) addClient(client server.Client) {
	m.clientMu.Lock()
	m.clientList.set(client.ClientOptions().ClientID, client)
	m.clientMu.Unlock()
}

func (m *monitor) setClientDisconnected(clientID string) {
	m.clientMu.Lock()
	defer m.clientMu.Unlock()


}

// deleteClient
func (m *monitor) deleteClient(clientID string) {
	m.clientMu.Lock()
	m.clientList.remove(clientID)
	m.clientMu.Unlock()
}

// deleteClientSubscriptions
func (m *monitor) deleteClientSubscriptions(clientID string) {
	m.subMu.Lock()
	defer m.subMu.Unlock()
	delete(m.subscriptions, clientID)
}

// addSubscription
func (m *monitor) addSubscription(clientID string, subscription *gmqtt.Subscription) {
	m.subMu.Lock()
	defer m.subMu.Unlock()
	subInfo := &SubscriptionInfo{
		ClientID: clientID,
		Qos:      subscription.QoS,
		Name:     subscription.GetFullTopicName(),
		At:       time.Now(),
	}
	if _, ok := m.subscriptions[clientID]; !ok {
		m.subscriptions[clientID] = newQuickList()
	}
	m.subscriptions[clientID].set(subInfo.Name, subInfo)
}

// deleteSubscription
func (m *monitor) deleteSubscription(clientID string, topicName string) {
	m.subMu.Lock()
	defer m.subMu.Unlock()
	if _, ok := m.subscriptions[clientID]; ok {
		m.subscriptions[clientID].remove(topicName)
	}
}

// GetClients ...
func (m *monitor) GetClients(offset, n int) ([]*ClientInfo, int, error) {
	rs := make([]*ClientInfo, 0)
	fn := func(elem *list.Element) {
		rs = append(rs, newClientInfo(elem.Value.(server.Client)))
	}
	m.clientMu.Lock()
	m.clientList.iterate(fn, offset, n)
	total := m.clientList.rows.Len()
	m.clientMu.Unlock()
	return rs, total, nil
}

func (m *monitor) getClientByID(clientID string) (server.Client, error) {
	if i, err := m.clientList.getByID(clientID); i != nil {
		return i.Value.(server.Client), nil
	} else {
		return nil, err
	}
}

// GetClientByID ...
func (m *monitor) GetClientByID(clientID string) (*ClientInfo, error) {
	m.clientMu.Lock()
	client, err := m.getClientByID(clientID)
	m.clientMu.Unlock()
	if err != nil {
		return nil, err
	}
	return newClientInfo(client), err
}

func statusText(client server.Client) string {
	if client.SessionInfo().IsExpired(time.Now()) {
		return Online
	} else {
		return Offline
	}
}

func (m *monitor) newSessionInfo(client server.Client, c config.Config) *SessionInfo {
	optsReader := client.ClientOptions()
	stats, _ := m.statsReader.GetClientStats(client.ClientOptions().ClientID)
	subStats, _ := m.subStatsReader.GetClientStats(optsReader.ClientID)
	rs := &SessionInfo{
		ClientID: optsReader.ClientID,
		Status:   statusText(client),
		//CleanSession:          optsReader.CleanSession(),
		Subscriptions: subStats.SubscriptionsCurrent,
		MaxInflight:   c.MQTT.MaxInflight,
		InflightLen:   stats.MessageStats.InflightCurrent,
		MaxMsgQueue:   c.MQTT.MaxQueuedMsg,
		MsgQueueLen:   stats.MessageStats.QueuedCurrent,
		//AwaitRelLen:           stats.AwaitRelCurrent,
		Qos0MsgDroppedTotal:   stats.MessageStats.Qos0.DroppedTotal.QueueFull,
		Qos1MsgDroppedTotal:   stats.MessageStats.Qos1.DroppedTotal.QueueFull,
		Qos2MsgDroppedTotal:   stats.MessageStats.Qos2.DroppedTotal.QueueFull,
		Qos0MsgDeliveredTotal: stats.MessageStats.Qos0.SentTotal,
		Qos1MsgDeliveredTotal: stats.MessageStats.Qos1.SentTotal,
		Qos2MsgDeliveredTotal: stats.MessageStats.Qos2.SentTotal,
		ConnectedAt:           common.Time(client.ConnectedAt()),
		//DisconnectedAt:        client.DisconnectedAt(),
	}
	return rs
}

// GetSessions ...
func (m *monitor) GetSessions(offset, n int) ([]*SessionInfo, int, error) {
	rs := make([]*SessionInfo, 0)
	fn := func(elem *list.Element) {
		rs = append(rs, m.newSessionInfo(elem.Value.(server.Client), m.config))
	}
	m.clientMu.Lock()
	m.clientList.iterate(fn, offset, n)
	total := m.clientList.rows.Len()
	m.clientMu.Unlock()
	return rs, total, nil
}

// GetSessionByID ...
func (m *monitor) GetSessionByID(clientID string) (*SessionInfo, error) {
	m.clientMu.Lock()
	client, err := m.getClientByID(clientID)
	m.clientMu.Unlock()
	if err != nil {
		return nil, err
	}
	return m.newSessionInfo(client, m.config), err
}

// GetClientSubscriptions ...
func (m *monitor) GetClientSubscriptions(clientID string, offset, n int) ([]*SubscriptionInfo, int, error) {
	m.subMu.Lock()
	defer m.subMu.Unlock()
	rs := make([]*SubscriptionInfo, 0)
	var err error
	var total int
	if _, ok := m.subscriptions[clientID]; ok {
		fn := func(elem *list.Element) {
			rs = append(rs, elem.Value.(*SubscriptionInfo))
		}
		total = m.subscriptions[clientID].rows.Len()
		err = m.subscriptions[clientID].iterate(fn, offset, n)
	}
	return rs, total, err
}

// SearchTopic ...
func (m *monitor) SearchTopic(query string) (rs []*SubscriptionInfo, err error) {
	m.subMu.Lock()
	defer m.subMu.Unlock()
	rs = make([]*SubscriptionInfo, 0)
	var info *SubscriptionInfo
	for _, sub := range m.subscriptions {
		fn := func(elem *list.Element) {
			info = elem.Value.(*SubscriptionInfo)
			if !strings.Contains(info.Name, query) {
				return
			}
			rs = append(rs, info)
		}
		err = sub.iterate(fn, 0, 999)
	}
	return
}
