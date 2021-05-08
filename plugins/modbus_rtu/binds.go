package modbus_rtu

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/e154/smart-home/plugins/node"
	"github.com/e154/smart-home/system/event_bus"
	"strings"
	"time"
)

type modbusRtu func(f string, address, count uint16, command []uint16) (result ModBusResponse)

func NewModbusRtu(eventBus event_bus.EventBus, actor *EntityActor) (modbus modbusRtu) {

	return func(f string, address, count uint16, command []uint16) (result ModBusResponse) {
		fmt.Printf("send message^ func(%s), address(%d), count(%d), comman(%b) \n", f, address, count, command)

		// time metric
		startTime := time.Now()

		// set callback func
		ch := make(chan node.MessageResponse)
		fn := func(topic string, msg node.MessageResponse) {
			items := strings.Split(topic, "/")
			entityId := items[len(items)-1]
			if entityId != actor.Id.String() {
				return
			}
			ch <- msg
		}
		var topic = actor.localTopic(fmt.Sprintf("resp/%s", actor.Id))
		eventBus.Subscribe(topic, fn)
		defer func() {
			eventBus.Unsubscribe(topic, fn)
		}()

		// send message
		actor.AttrMu.Lock()
		attrsSerial := actor.Attrs.Serialize()
		actor.AttrMu.Unlock()

		properties, err := json.Marshal(attrsSerial)
		if err != nil {
			log.Error(err.Error())
			result.Error = err.Error()
			return
		}

		cmd := ModBusCommand{
			Function: f,
			Address:  address,
			Count:    count,
			Command:  command,
		}
		b, err := json.Marshal(cmd)
		if err != nil {
			result.Error = err.Error()
			log.Error(err.Error())
			return
		}
		msg := node.MessageRequest{
			EntityId:   actor.Id,
			DeviceType: DeviceTypeModbusRtu,
			Properties: properties,
			Command:    b,
		}

		eventBus.Publish(actor.localTopic(fmt.Sprintf("req/%s", actor.Id)), msg)

		// wait response
		ticker := time.NewTimer(time.Second * 5)
		defer ticker.Stop()

		var done bool
		for {
			if done {
				break
			}
			select {
			case <-ticker.C:
				err = errors.New("timeout")
				done = true
			case resp := <-ch:
				if resp.EntityId != actor.Id {
					continue
				}
				if err = json.Unmarshal(resp.Response, &result); err != nil {
					log.Error(err.Error())
				}
				done = true
			}
		}

		result.Time = time.Since(startTime).Seconds()

		return
	}
}
