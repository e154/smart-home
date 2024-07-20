// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2024, Filippov Alex
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

package ble

import (
	"go.uber.org/atomic"
	"sync"
	"tinygo.org/x/bluetooth"
)

type Ble struct {
	isScan      *atomic.Bool
	scanAddress *bluetooth.UUID
	devMX       sync.Mutex
	device      *bluetooth.Device
	connected   *atomic.Bool
	adapter     *bluetooth.Adapter
	address     string
	timeout     int64
}

func NewBle() *Ble {
	ble := &Ble{
		isScan:    atomic.NewBool(false),
		connected: atomic.NewBool(false),
		adapter:   bluetooth.DefaultAdapter,
	}

	ble.adapter.SetConnectHandler(func(device bluetooth.Device, connected bool) {
		log.Infof("bluetooth device: %s, connected: %t", device.Address.String(), connected)
		ble.connected.Store(connected)
		ble.device = &device

		//if connected || ble.address == "" {
		//	return
		//}
		//
		//time.Sleep(time.Second * 10)
		//ble.Connect(ble.address, ble.timeout)
	})

	return ble
}

func (b *Ble) Disconnect() error {
	if !b.connected.Load() || b.device == nil {
		return nil
	}

	b.connected.Store(false)
	b.address = ""
	if err := b.device.Disconnect(); err != nil {
		return err
	}
	b.device = nil
	return nil
}

func (b *Ble) GetServices(address string, timeout int64) ([]bluetooth.DeviceService, error) {
	device, err := b.Connect(address, timeout)
	if err != nil {
		return nil, err
	}

	// Get a list of services
	return device.DiscoverServices(nil)
}

func (b *Ble) GetCharacteristics(address string, chars []bluetooth.UUID, timeout int64) ([]bluetooth.DeviceCharacteristic, error) {

	var characteristic = []bluetooth.DeviceCharacteristic{}

	services, err := b.GetServices(address, timeout)
	if err != nil {
		return nil, err
	}

	for _, service := range services {
		discoverCharacteristics, _ := service.DiscoverCharacteristics(chars)
		if len(discoverCharacteristics) == 0 {
			continue
		}
		characteristic = append(characteristic, discoverCharacteristics...)
	}

	return characteristic, nil
}

type Cache struct {
	sync.Mutex
	pull map[bluetooth.UUID]map[bluetooth.UUID]struct{}
}

func NewCache() *Cache {
	return &Cache{
		pull: make(map[bluetooth.UUID]map[bluetooth.UUID]struct{}),
	}
}

func (c *Cache) Get(key bluetooth.UUID) (bluetooth.UUID, bool) {
	c.Lock()
	defer c.Unlock()
	for k, values := range c.pull {
		if _, ok := values[key]; ok {
			return k, ok
		}
	}
	return bluetooth.UUID{}, false
}

func (c *Cache) Put(key, value bluetooth.UUID) {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.pull[key]; !ok {
		c.pull[key] = make(map[bluetooth.UUID]struct{})
	}
	c.pull[key][value] = struct{}{}
}
