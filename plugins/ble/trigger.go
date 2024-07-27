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
	"github.com/e154/bus"
	"github.com/e154/smart-home/plugins/triggers"
	"reflect"
	"sync"
	"tinygo.org/x/bluetooth"
)

var _ triggers.ITrigger = (*Trigger)(nil)

type Trigger struct {
	eventBus     bus.Bus
	msgQueue     bus.Bus
	functionName string
	name         string
	ble          map[string]map[string]*Ble
}

func NewTrigger(eventBus bus.Bus) triggers.ITrigger {
	return &Trigger{
		eventBus:     eventBus,
		msgQueue:     bus.NewBus(),
		functionName: FunctionName,
		name:         Name,
		ble:          make(map[string]map[string]*Ble),
	}
}

func (t *Trigger) Name() string {
	return t.name
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

	if _, ok := t.ble[address.String()][characteristic.String()]; ok {
		return fmt.Errorf("a trigger with such parameters already exists")
	}

	var timeout, connectionTimeout int64 = 5, 5
	ble := NewBle(address.String(), timeout, connectionTimeout, false)

	char, err := bluetooth.ParseUUID(characteristic.String())
	if err != nil {
		return err
	}

	callback := reflect.ValueOf(options.Handler)
	err = ble.Subscribe(char, func(bytes []byte) {
		callback.Call([]reflect.Value{reflect.ValueOf(""), reflect.ValueOf(bytes)})
	})

	if err != nil {
		return err
	}

	if _, ok := t.ble[address.String()]; !ok {
		t.ble[address.String()] = make(map[string]*Ble)
	}
	t.ble[address.String()][characteristic.String()] = ble

	log.Infof("trigger '%s' subscribed to '%s'", t.name, characteristic.String())

	return nil
}

// Unsubscribe ...
func (t *Trigger) Unsubscribe(options triggers.Subscriber) error {

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

	ble, ok := t.ble[address.String()][characteristic.String()]
	if !ok {
		return nil
	}

	_ = ble.Disconnect()

	delete(t.ble, address.String())

	log.Infof("unsubscribe from %s", characteristic.String())

	return nil
}

// FunctionName ...
func (t *Trigger) FunctionName() string {
	return t.functionName
}
