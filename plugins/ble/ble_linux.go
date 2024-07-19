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
	device      *bluetooth.Device
}

func NewBle() *Ble {
	return &Ble{
		isScan: atomic.NewBool(false),
	}
}

func (b *Ble) Disconnect() error {
	if b.device == nil {
		return nil
	}

	if err := b.device.Disconnect(); err != nil {
		return err
	}
	b.device = nil
	return nil
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

		if b.scanAddress == nil || device.Address.String() == "" {
			log.Infof("found device: %s, RSSI: %v, LocalName: %s, payload: %v", scanResult.Address.String(), scanResult.RSSI, scanResult.LocalName(), scanResult.AdvertisementPayload)
		}

		if b.scanAddress == nil || device.Address.String() != b.scanAddress.String() {
			return
		}

		log.Infof("found device: %s, RSSI: %v, LocalName: %s, payload: %v", scanResult.Address.String(), scanResult.RSSI, scanResult.LocalName(), scanResult.AdvertisementPayload)

		adapter.StopScan()

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

func (b *Ble) connectBluetooth(address string, timeout int64) (*bluetooth.Device, error) {

	if b.device != nil {
		return b.device, nil
	}

	adapter := bluetooth.DefaultAdapter
	_ = adapter.Enable()

	if timeout == 0 {
		timeout = 1
	}

	mac, err := bluetooth.ParseMAC(address)
	if err != nil {
		return nil, err
	}
	device, err := adapter.Connect(bluetooth.Address{
		MACAddress: bluetooth.MACAddress{
			MAC: mac,
		},
	}, bluetooth.ConnectionParams{
		ConnectionTimeout: bluetooth.NewDuration(time.Second * time.Duration(timeout)),
		Timeout:           bluetooth.NewDuration(time.Second * time.Duration(timeout)),
	})
	if err != nil {
		return nil, err
	}

	//log.Infof("connected: %s", device.Address.String())

	b.device = &device

	return &device, nil
}

func (b *Ble) Write(address string, char bluetooth.UUID, timeout int64, request []byte, withResponse bool) ([]byte, error) {

	device, err := b.connectBluetooth(address, timeout)
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

			log.Infof("write: %x --> %s", request, address)

			if _, err = characteristic.WriteWithoutResponse(request); err != nil {
				return nil, err
			}
			return []byte{}, nil
		}
	}

	return nil, nil
}

func (b *Ble) Read(address string, char bluetooth.UUID, timeout int64) ([]byte, error) {

	device, err := b.connectBluetooth(address, timeout)
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

			payload := make([]byte, 1024)
			i, err := characteristic.Read(payload)
			if err != nil {
				return nil, err
			}
			return payload[:uint32(i)], nil
		}
	}

	return nil, nil
}

func (b *Ble) Subscribe(address string, char bluetooth.UUID, timeout int64, handler func([]byte)) error {

	device, err := b.connectBluetooth(address, timeout)
	if err != nil {
		return err
	}

	// Get a list of services
	services, err := device.DiscoverServices(nil)
	if err != nil {
		return err
	}

	// Iterate services
	for _, service := range services {
		// Get a list of characteristics below the service
		characteristics, err := service.DiscoverCharacteristics(nil)
		if err != nil {
			return err
		}

		// Iterate characteristics
		for _, characteristic := range characteristics {

			if characteristic.UUID() != char {
				continue
			}

			return characteristic.EnableNotifications(handler)
		}
	}

	return nil
}
