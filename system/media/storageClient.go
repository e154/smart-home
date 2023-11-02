// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

package media

import (
	"time"

	"github.com/google/uuid"

	"github.com/deepch/vdk/av"
)

// ClientAdd Add New Client to Translations
func (obj *StorageST) ClientAdd(streamID string, channelID string, mode int) (string, chan *av.Packet, chan *[]byte, error) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	streamTmp, ok := obj.Streams[streamID]
	if !ok {
		return "", nil, nil, ErrorStreamNotFound
	}
	//Generate UUID client
	cid := uuid.New().String()
	chAV := make(chan *av.Packet, 2000)
	chRTP := make(chan *[]byte, 2000)
	channelTmp, ok := streamTmp.Channels[channelID]
	if !ok {
		return "", nil, nil, ErrorStreamNotFound
	}

	channelTmp.clients[cid] = ClientST{mode: mode, outgoingAVPacket: chAV, outgoingRTPPacket: chRTP, signals: make(chan int, 100)}
	channelTmp.ack = time.Now()
	streamTmp.Channels[channelID] = channelTmp
	obj.Streams[streamID] = streamTmp
	return cid, chAV, chRTP, nil

}

// ClientDelete Delete Client
func (obj *StorageST) ClientDelete(streamID string, cid string, channelID string) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if _, ok := obj.Streams[streamID]; ok {
		delete(obj.Streams[streamID].Channels[channelID].clients, cid)
	}
}

// ClientHas check is client ext
func (obj *StorageST) ClientHas(streamID string, channelID string) bool {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	streamTmp, ok := obj.Streams[streamID]
	if !ok {
		return false
	}
	channelTmp, ok := streamTmp.Channels[channelID]
	if !ok {
		return false
	}
	if time.Now().Sub(channelTmp.ack).Seconds() > 30 {
		return false
	}
	return true
}
