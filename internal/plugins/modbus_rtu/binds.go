// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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

package modbus_rtu

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/e154/smart-home/internal/plugins/node"
	"github.com/e154/smart-home/pkg/apperr"

	"github.com/e154/bus"
	"github.com/pkg/errors"
	"go.uber.org/atomic"
)

type modbusRtu func(f string, address, count uint16, command []uint16) (result ModBusResponse)

// NewModbusRtu ...
func NewModbusRtu(eventBus bus.Bus, actor *Actor) (modbus modbusRtu) {

	var isStarted = atomic.NewBool(false)

	return func(f string, address, count uint16, command []uint16) (result ModBusResponse) {
		//log.Debugf("send message func(%s), address(%d), count(%d), command(%b)", f, address, count, command)

		var err error
		defer func() {
			if err != nil {
				result.Error = err.Error()
			}
		}()

		if isStarted.Load() {
			err = fmt.Errorf("in process")
			return
		}
		isStarted.Store(true)
		defer isStarted.Store(false)

		// time metric
		startTime := time.Now()

		// set callback func
		ch := make(chan node.MessageResponse)
		defer close(ch)
		fn := func(topic string, msg node.MessageResponse) {
			items := strings.Split(topic, "/")
			entityId := items[len(items)-1]
			if entityId != actor.Id.String() {
				return
			}
			ch <- msg
		}
		var topic = actor.localTopic(fmt.Sprintf("resp/%s", actor.Id))
		if err = eventBus.Subscribe(topic, fn, false); err != nil {
			log.Error(err.Error())
			return
		}
		defer func() {
			_ = eventBus.Unsubscribe(topic, fn)
		}()

		var properties []byte
		if properties, err = json.Marshal(actor.Settings().Serialize()); err != nil {
			log.Error(err.Error())
			return
		}

		cmd := ModBusCommand{
			Function: f,
			Address:  address,
			Count:    count,
			Command:  command,
		}
		var requestCommand []byte
		if requestCommand, err = json.Marshal(cmd); err != nil {
			log.Error(err.Error())
			return
		}
		msg := node.MessageRequest{
			EntityId:   actor.Id,
			DeviceType: DeviceTypeModbusRtu,
			Properties: properties,
			Command:    requestCommand,
		}

		eventBus.Publish(actor.localTopic(fmt.Sprintf("req/%s", actor.Id)), msg)

		// wait response
		ticker := time.NewTimer(time.Second * 1)
		defer ticker.Stop()

		select {
		case <-ticker.C:
			err = errors.Wrap(apperr.ErrTimeout, "wait timeout")
		case resp := <-ch:
			if err = json.Unmarshal(resp.Response, &result); err != nil {
				log.Error(err.Error())
			}
		}

		result.Time = time.Since(startTime).Seconds()

		return
	}
}
