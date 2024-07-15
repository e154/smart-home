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
	"time"

	"go.uber.org/atomic"
	"tinygo.org/x/bluetooth"
)

type Ble struct {
	isScan      *atomic.Bool
	scanAddress *bluetooth.UUID
}

func NewBle() *Ble {
	return &Ble{
		isScan: atomic.NewBool(false),
	}
}

func (b *Ble) Scan(address *bluetooth.UUID) {
	if !b.isScan.CompareAndSwap(false, true) {
		return
	}
	log.Info("Start scan")
	b.scanAddress = address

	defer func() {
		log.Info("Stop scan")
		b.isScan.Store(false)
	}()

	adapter := bluetooth.DefaultAdapter
	_ = adapter.Enable()

	go func() {
		// Start scanning and define callback for scan results
		if err := adapter.Scan(b.onScan); err != nil {
			log.Error(err.Error())
		}
	}()

	select {
	case <-time.After(time.Second * 10):
		_ = adapter.StopScan()
	}
}

func (b *Ble) onScan(adapter *bluetooth.Adapter, scanResult bluetooth.ScanResult) {

	// Start connecting in a goroutine to not block
	go func() {
		device, err := adapter.Connect(scanResult.Address, bluetooth.ConnectionParams{})
		if err != nil {
			return
		}

		if b.scanAddress != nil && device.Address.String() != b.scanAddress.String() {
			return
		}

		log.Infof("found device: %s, RSSI: %v, LocalName: %s, payload: %v", scanResult.Address.String(), scanResult.RSSI, scanResult.LocalName(), scanResult.AdvertisementPayload)

		// Call connect callback
		b.onScanConnect(device)
	}()
}

func (b *Ble) onScanConnect(device bluetooth.Device) {

	defer func() {
		if r := recover(); r != nil {
			log.Warn("Recovered")
		}
	}()

	// Get a list of services
	services, err := device.DiscoverServices(nil)
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Iterate services
	for _, service := range services {

		log.Infof("service: %s", service.UUID().String())

		// Get a list of characteristics below the service
		characteristics, err := service.DiscoverCharacteristics(nil)
		if err != nil {
			log.Error(err.Error())
			return
		}

		// Iterate characteristics
		for _, char := range characteristics {
			log.Infof("characteristic: %s", char.UUID().String())
		}
	}

}

func (b *Ble) connectBluetooth(address bluetooth.UUID) (*bluetooth.Device, error) {

	adapter := bluetooth.DefaultAdapter
	_ = adapter.Enable()

	device, err := adapter.Connect(bluetooth.Address{UUID: address}, bluetooth.ConnectionParams{
		ConnectionTimeout: bluetooth.NewDuration(time.Second * 2),
		Timeout:           bluetooth.NewDuration(time.Second * 2),
	})
	if err != nil {
		return nil, err
	}

	//log.Infof("connected: %s", device.Address.String())

	return &device, nil
}

func (b *Ble) Write(address, char bluetooth.UUID, payload []byte) (int, error) {

	device, err := b.connectBluetooth(address)
	if err != nil {
		return 0, err
	}

	// Get a list of services
	services, err := device.DiscoverServices(nil)
	if err != nil {
		return 0, err
	}

	// Iterate services
	for _, service := range services {
		// Get a list of characteristics below the service
		characteristics, err := service.DiscoverCharacteristics(nil)
		if err != nil {
			return 0, err
		}

		// Iterate characteristics
		for _, characteristic := range characteristics {

			if characteristic.UUID() != char {
				continue
			}

			log.Infof("write: %x --> %s", payload, address)
			n, err := characteristic.Write(payload)
			if err != nil {
				return 0, err
			}
			return n, nil
		}
	}

	return 0, nil
}

func (b *Ble) Read(address, char bluetooth.UUID) ([]byte, error) {

	device, err := b.connectBluetooth(address)
	if err != nil {
		return nil, err
	}

	// Get a list of services
	services, err := device.DiscoverServices(nil)
	if err != nil {
		return nil, err
	}

	// Iterate services
	for _, service := range services {
		// Get a list of characteristics below the service
		characteristics, err := service.DiscoverCharacteristics(nil)
		if err != nil {
			return nil, err
		}

		// Iterate characteristics
		for _, characteristic := range characteristics {

			if characteristic.UUID() != char {
				continue
			}

			payload := make([]byte, 255)
			if _, err = characteristic.Read(payload); err != nil {
				return nil, err
			}
			return payload, nil
		}
	}

	return nil, nil
}
