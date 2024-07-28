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

// EventStateChanged

package ble

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/e154/bus"
	"tinygo.org/x/bluetooth"

	"github.com/e154/smart-home/plugins/triggers"
)

var _ triggers.ITrigger = (*Trigger)(nil)

type Trigger struct {
	eventBus     bus.Bus
	msgQueue     bus.Bus
	functionName string
	name         string
	ticker       *time.Ticker
	sync.Mutex
	devices map[string]map[string]*TriggerParams
}

func NewTrigger(eventBus bus.Bus) *Trigger {
	trigger := &Trigger{
		eventBus:     eventBus,
		msgQueue:     bus.NewBus(),
		functionName: FunctionName,
		name:         Name,
		devices:      make(map[string]map[string]*TriggerParams),
	}

	go func() {
		const pause = 10
		trigger.ticker = time.NewTicker(time.Second * time.Duration(pause))

		for range trigger.ticker.C {
			trigger.Lock()
			for _, device := range trigger.devices {
				for _, characteristic := range device {
					if !characteristic.connected.Load() {
						trigger.Connect(characteristic, false)
					}
				}
			}
			trigger.Unlock()
		}
	}()

	return trigger
}

func (t *Trigger) Name() string {
	return t.name
}

func (t *Trigger) Shutdown() {
	if t.ticker != nil {
		t.ticker.Stop()
		t.ticker = nil
	}

	t.Lock()
	defer t.Unlock()

	for _, device := range t.devices {
		for _, characteristic := range device {
			characteristic.Disconnect()
		}
	}
}

func (t *Trigger) AsyncAttach(wg *sync.WaitGroup) {

	wg.Done()
}

// Subscribe ...
func (t *Trigger) Subscribe(options triggers.Subscriber) error {

	if options.Payload == nil {
		return fmt.Errorf("payload is nil")
	}

	address, ok := options.Payload[AttrAddress]
	if !ok {
		return fmt.Errorf("address attribute is nil")
	}

	characteristic, ok := options.Payload[AttrCharacteristic]
	if !ok {
		return fmt.Errorf("characteristic attribute is nil")
	}

	if _, ok := t.devices[address.String()][characteristic.String()]; ok {
		return fmt.Errorf("a trigger with such parameters already exists")
	}

	t.Lock()
	defer t.Unlock()

	err := t.Connect(&TriggerParams{options: options}, true)

	return err
}

// Unsubscribe ...
func (t *Trigger) Unsubscribe(options triggers.Subscriber) error {

	t.Lock()
	defer t.Unlock()

	if options.Payload == nil {
		return fmt.Errorf("payload is nil")
	}

	address, ok := options.Payload[AttrAddress]
	if !ok {
		return fmt.Errorf("address attribute is nil")
	}

	characteristic, ok := options.Payload[AttrCharacteristic]
	if !ok {
		return fmt.Errorf("characteristic attribute is nil")
	}

	device, ok := t.devices[address.String()][characteristic.String()]
	if !ok {
		return nil
	}

	if device != nil {
		_ = device.Disconnect()
	}

	delete(t.devices, address.String())

	log.Infof("unsubscribe from %s", characteristic.String())

	return nil
}

// FunctionName ...
func (t *Trigger) FunctionName() string {
	return t.functionName
}

func (t *Trigger) Connect(params *TriggerParams, firstTime bool) error {

	options := params.options
	address := options.Payload[AttrAddress].String()
	characteristic := options.Payload[AttrCharacteristic].String()

	var timeout, connectionTimeout int64 = 5, 5
	params.Ble = NewBle(address, timeout, connectionTimeout, false)

	char, err := bluetooth.ParseUUID(characteristic)
	if err != nil {
		return err
	}

	callback := reflect.ValueOf(options.Handler)
	err = params.Subscribe(char, func(bytes []byte) {
		callback.Call([]reflect.Value{reflect.ValueOf(""), reflect.ValueOf(bytes)})
	})

	if err != nil && firstTime {
		log.Infof("trigger '%s' is not subscribed to '%s' but we are trying", t.name, characteristic)
	}

	if _, ok := t.devices[address]; !ok {
		t.devices[address] = make(map[string]*TriggerParams)
	}
	t.devices[address][characteristic] = params

	if err == nil && firstTime {
		log.Infof("trigger '%s' subscribed to '%s'", t.name, characteristic)
	}

	return nil
}
