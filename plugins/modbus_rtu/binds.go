package modbus_rtu

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/e154/smart-home/plugins/node"
	"github.com/e154/smart-home/system/event_bus"
	"go.uber.org/atomic"
	"strings"
	"time"
)

type modbusRtu func(f string, address, count uint16, command []uint16) (result ModBusResponse)

func NewModbusRtu(eventBus event_bus.EventBus, actor *EntityActor) (modbus modbusRtu) {

	var isStarted = atomic.NewBool(false)

	return func(f string, address, count uint16, command []uint16) (result ModBusResponse) {
		//fmt.Printf("send message^ func(%s), address(%d), count(%d), command(%b) \n", f, address, count, command)

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
			eventBus.Unsubscribe(topic, fn)
		}()

		// send message
		actor.AttrMu.Lock()
		attrsSerial := actor.Attrs.Serialize()
		actor.AttrMu.Unlock()

		var properties []byte
		if properties, err = json.Marshal(attrsSerial); err != nil {
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
			err = errors.New("timeout")
		case resp := <-ch:
			if err = json.Unmarshal(resp.Response, &result); err != nil {
				log.Error(err.Error())
			}
		}

		result.Time = time.Since(startTime).Seconds()

		return
	}
}
