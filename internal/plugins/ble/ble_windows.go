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

	"tinygo.org/x/bluetooth"
)

func (b *Ble) scan(address *bluetooth.UUID) {
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

func (b *Ble) connect() (*bluetooth.Device, error) {

	if b.debug {
		log.Debugf("try connecting to %v", b.address)
	}

	b.devMX.Lock()
	defer b.devMX.Unlock()

	if b.connected.Load() {
		if b.debug {
			log.Debugf("there is still a connection %v", b.address)
		}
		return b.device, nil
	}

	mac, err := bluetooth.ParseMAC(b.address)
	if err != nil {
		return nil, err
	}

	if b.debug {
		log.Debugf("try to turn on the device")
		if err := b.adapter.Enable(); err != nil {
			log.Error(err.Error())
		}
	} else {
		b.adapter.Enable()
	}

	if b.debug {
		log.Debugf("Connect starts a connection attempt to %s", b.address)
	}
	device, err := b.adapter.Connect(bluetooth.Address{
		MACAddress: bluetooth.MACAddress{
			MAC: mac,
		},
	}, bluetooth.ConnectionParams{
		ConnectionTimeout: bluetooth.NewDuration(time.Second * time.Duration(b.connectionTimeout)),
		Timeout:           bluetooth.NewDuration(time.Second * time.Duration(b.timeout)),
	})
	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (b *Ble) write(char bluetooth.UUID, request []byte, withResponse bool) ([]byte, error) {

	characteristics, err := b.GetCharacteristics([]bluetooth.UUID{char})
	if err != nil {
		return nil, err
	}

	// Iterate characteristics
	for _, characteristic := range characteristics {

		if characteristic.UUID() != char {
			continue
		}

		log.Infof("write: %x --> %s", request, b.address)

		if withResponse {
			payload := make([]byte, len(request))
			copy(payload, request)
			i, err := characteristic.Write(payload)
			if err != nil {
				return nil, err
			}
			log.Infof("read: %x <-- %s", payload[:uint32(i)], b.address)
			return payload[:uint32(i)], nil
		} else {
			if _, err = characteristic.WriteWithoutResponse(request); err != nil {
				return nil, err
			}
			return []byte{}, nil
		}
	}

	return nil, nil
}

func (b *Ble) read(char bluetooth.UUID) ([]byte, error) {

	characteristics, err := b.GetCharacteristics([]bluetooth.UUID{char})
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
		log.Infof("read: %x <-- %s", payload[:uint32(i)], b.address)
		return payload[:uint32(i)], nil
	}

	return nil, nil
}

func (b *Ble) subscribe(char bluetooth.UUID, handler func([]byte)) error {

	characteristics, err := b.GetCharacteristics([]bluetooth.UUID{char})
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

	return nil
}
