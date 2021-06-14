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
	"errors"
	"github.com/DrmagicE/gmqtt"
	"github.com/DrmagicE/gmqtt/config"
	"github.com/DrmagicE/gmqtt/persistence/subscription"
	"github.com/DrmagicE/gmqtt/server"
	"github.com/e154/smart-home/common"
	"strings"
	"sync"
	"time"
)

type store struct {
	clientMu       sync.RWMutex
	clientIndexer  *Indexer
	subMu          sync.RWMutex
	subIndexer     *Indexer
	config         config.Config
	statsReader    server.StatsReader
	subStatsReader subscription.StatsReader
	clientService  server.ClientService
}

func newStore(statsReader server.StatsReader,
	subStatsReader subscription.StatsReader,
	clientService server.ClientService) *store {
	return &store{
		clientIndexer:  NewIndexer(),
		subIndexer:     NewIndexer(),
		statsReader:    statsReader,
		subStatsReader: subStatsReader,
		clientService:  clientService,
	}
}

func (s *store) addSubscription(clientID string, sub *gmqtt.Subscription) {
	s.subMu.Lock()
	defer s.subMu.Unlock()

	subInfo := &SubscriptionInfo{
		TopicName:         sub.GetFullTopicName(),
		Id:                sub.ID,
		Qos:               uint32(sub.QoS),
		NoLocal:           sub.NoLocal,
		RetainAsPublished: sub.RetainAsPublished,
		RetainHandling:    uint32(sub.RetainHandling),
		ClientID:          clientID,
	}
	key := clientID + "_" + sub.GetFullTopicName()
	s.subIndexer.Set(key, subInfo)

}

func (s *store) removeSubscription(clientID string, topicName string) {
	s.subMu.Lock()
	defer s.subMu.Unlock()
	s.subIndexer.Remove(clientID + "_" + topicName)
}

func (s *store) addClient(client server.Client) {
	c := newClientInfo(client)
	s.clientMu.Lock()
	s.clientIndexer.Set(c.ClientID, c)
	s.clientMu.Unlock()
}

func (s *store) setClientDisconnected(clientID string) {
	s.clientMu.Lock()
	defer s.clientMu.Unlock()
	l := s.clientIndexer.GetByID(clientID)
	if l == nil {
		return
	}
	l.Value.(*ClientInfo).DisconnectedAt = common.Time(time.Now())
}

func (s *store) removeClient(clientID string) {
	s.clientMu.Lock()
	s.clientIndexer.Remove(clientID)
	s.clientMu.Unlock()
}

// GetClientByID returns the client information for the given client id.
func (s *store) GetClientByID(clientID string) *ClientInfo {
	s.clientMu.RLock()
	defer s.clientMu.RUnlock()
	c := s.getClientByIDLocked(clientID)
	fillClientInfo(c, s.statsReader)
	return c
}

func (s *store) getClientByIDLocked(clientID string) *ClientInfo {
	if i := s.clientIndexer.GetByID(clientID); i != nil {
		return i.Value.(*ClientInfo)
	}
	return nil
}

func fillClientInfo(c *ClientInfo, stsReader server.StatsReader) {
	if c == nil {
		return
	}
	sts, ok := stsReader.GetClientStats(c.ClientID)
	if !ok {
		return
	}
	c.SubscriptionsCurrent = uint32(sts.SubscriptionStats.SubscriptionsCurrent)
	c.SubscriptionsTotal = uint32(sts.SubscriptionStats.SubscriptionsTotal)
	c.PacketsReceivedBytes = sts.PacketStats.BytesReceived.Total
	c.PacketsReceivedNums = sts.PacketStats.ReceivedTotal.Total
	c.PacketsSendBytes = sts.PacketStats.BytesSent.Total
	c.PacketsSendNums = sts.PacketStats.SentTotal.Total
	c.MessageDropped = sts.MessageStats.GetDroppedTotal()
	c.InflightLen = uint32(sts.MessageStats.InflightCurrent)
	c.QueueLen = uint32(sts.MessageStats.QueuedCurrent)
}

// GetClients
func (s *store) GetClients(limit, offset uint) (rs []*ClientInfo, total uint32, err error) {
	rs = make([]*ClientInfo, 0)
	fn := func(elem *list.Element) {
		c := elem.Value.(*ClientInfo)
		fillClientInfo(c, s.statsReader)
		rs = append(rs, elem.Value.(*ClientInfo))
	}
	s.clientMu.RLock()
	defer s.clientMu.RUnlock()
	s.clientIndexer.Iterate(fn, offset, limit)
	return rs, uint32(s.clientIndexer.Len()), nil
}

func (s *store) newSessionInfo(client server.Client, c config.Config) *SessionInfo {
	optsReader := client.ClientOptions()
	stats, _ := s.statsReader.GetClientStats(client.ClientOptions().ClientID)
	subStats, _ := s.subStatsReader.GetClientStats(optsReader.ClientID)
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

// GetSessionByID ...
func (s *store) GetSessionByID(clientID string) (*SessionInfo, error) {
	client := s.GetClientByID(clientID)
	if client == nil {
		return nil, errors.New("not found")
	}
	return s.newSessionInfo(s.clientService.GetClient(clientID), s.config), nil
}

// GetClientSubscriptions ...
func (s *store) GetClientSubscriptions(clientID string, offset, n uint) ([]*SubscriptionInfo, int, error) {
	s.subMu.Lock()
	defer s.subMu.Unlock()
	rs := make([]*SubscriptionInfo, 0)
	var err error
	var total int
	fn := func(elem *list.Element) {
		info := elem.Value.(*SubscriptionInfo)
		if info.ClientID != clientID {
			return
		}
		rs = append(rs, info)
	}
	s.subIndexer.Iterate(fn, offset, n)
	total = s.subIndexer.Len()
	return rs, total, err
}

// SearchTopic ...
func (s *store) SearchTopic(query string) (rs []*SubscriptionInfo, err error) {
	s.subMu.Lock()
	defer s.subMu.Unlock()
	rs = make([]*SubscriptionInfo, 0)
	var info *SubscriptionInfo
	fn := func(elem *list.Element) {
		info = elem.Value.(*SubscriptionInfo)
		if !strings.Contains(info.TopicName, query) {
			return
		}
		rs = append(rs, info)
	}
	s.subIndexer.Iterate(fn, 0, 999)
	return
}

// GetSessions ...
func (s *store) GetSessions(offset, n uint) ([]*SessionInfo, int, error) {
	rs := make([]*SessionInfo, 0)
	fn := func(elem *list.Element) {
		c := elem.Value.(*ClientInfo)
		fillClientInfo(c, s.statsReader)
		rs = append(rs, s.newSessionInfo(s.clientService.GetClient(c.ClientID), s.config))
	}
	s.clientMu.RLock()
	defer s.clientMu.RUnlock()
	s.clientIndexer.Iterate(fn, offset, n)
	total := s.clientIndexer.Len()
	return rs, total, nil
}
