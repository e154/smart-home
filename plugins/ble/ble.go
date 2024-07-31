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
	isScan            *atomic.Bool
	scanAddress       *bluetooth.UUID
	devMX             sync.Mutex
	device            *bluetooth.Device
	connected         *atomic.Bool
	adapter           *bluetooth.Adapter
	address           string
	timeout           int64
	connectionTimeout int64
	debug             bool
}

func NewBle(address string, timeout, connectionTimeout int64, debug bool) *Ble {
	ble := &Ble{
		isScan:            atomic.NewBool(false),
		connected:         atomic.NewBool(false),
		adapter:           bluetooth.DefaultAdapter,
		timeout:           timeout,
		connectionTimeout: connectionTimeout,
		address:           address,
		debug:             debug,
	}

	ble.adapter.SetConnectHandler(func(device bluetooth.Device, connected bool) {
		log.Infof("bluetooth device: %s, connected: %t", device.Address.String(), connected)
		ble.connected.Store(connected)
		ble.device = &device
	})

	return ble
}

func (b *Ble) Connect() error {
	_, err := b.connect()
	return err
}

func (b *Ble) Disconnect() error {
	if !b.connected.Load() || b.device == nil {
		return nil
	}

	b.connected.Store(false)
	if err := b.device.Disconnect(); err != nil {
		return err
	}
	b.device = nil
	return nil
}

func (b *Ble) IsConnected() bool {
	return b.connected.Load()
}

func (b *Ble) Scan(param *string) {

	if param == nil {
		b.scan(nil)
		return
	}

	address, err := bluetooth.ParseUUID(*param)
	if err != nil {
		b.scan(nil)
	} else {
		b.scan(&address)
	}
}

func (b *Ble) Write(c string, request []byte, withResponse bool) ([]byte, error) {
	char, err := bluetooth.ParseUUID(c)
	if err != nil {
		return nil, err
	}
	return b.write(char, request, withResponse)
}

func (b *Ble) Read(c string) ([]byte, error) {
	char, err := bluetooth.ParseUUID(c)
	if err != nil {
		return nil, err
	}
	return b.read(char)
}

func (b *Ble) Subscribe(c string, handler func([]byte)) error {
	char, err := bluetooth.ParseUUID(c)
	if err != nil {
		return err
	}
	return b.subscribe(char, handler)
}

func (b *Ble) GetServices() ([]bluetooth.DeviceService, error) {
	device, err := b.connect()
	if err != nil {
		return nil, err
	}

	if b.debug {
		log.Debugf("device %v get services", b.address)
	}

	// Get a list of services
	return device.DiscoverServices(nil)
}

func (b *Ble) GetCharacteristics(chars []bluetooth.UUID) ([]bluetooth.DeviceCharacteristic, error) {

	if b.debug {
		log.Debugf("device %v get characteristics %v", b.address, chars)
	}

	var characteristic []bluetooth.DeviceCharacteristic

	services, err := b.GetServices()
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

	if b.debug {
		log.Debugf("device %v found %d characteristics", b.address, len(characteristic))
	}
	return characteristic, nil
}
