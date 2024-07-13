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
	"strconv"
	"tinygo.org/x/bluetooth"
)

type Ble struct {
}

func (b *Ble) Scan() {
	adapter := bluetooth.DefaultAdapter
	// Enable adapter
	err := adapter.Enable()
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Start scanning and define callback for scan results
	err = adapter.Scan(b.onScan)
	if err != nil {
		log.Error(err.Error())
	}
}

// var uid, _ = bluetooth.ParseUUID("bb917511-7cfa-c6a4-de37-935e8b9f22ea") //govee
// var uid, _ = bluetooth.ParseUUID("618cd365-792d-9ae0-279b-d8e033ef68a4") //sber
var uid, _ = bluetooth.ParseUUID("7f81005b-a978-ee00-d555-8a0b126eddc6") // Oclean

func (b *Ble) onScan(adapter *bluetooth.Adapter, scanResult bluetooth.ScanResult) {
	//log.Infof("found device: %s, RSSI: %v, LocalName: %s, payload: %v", scanResult.Address.String(), scanResult.RSSI, scanResult.LocalName(), scanResult.AdvertisementPayload)

	if scanResult.LocalName() == "Govee_H617A_3167" {
		log.Infof("found device: %s, RSSI: %v, LocalName: %s, payload: %v", scanResult.Address.String(), scanResult.RSSI, scanResult.LocalName(), scanResult.AdvertisementPayload)

		// Start connecting in a goroutine to not block
		go func() {
			device, err := adapter.Connect(scanResult.Address, bluetooth.ConnectionParams{})
			if err != nil {
				log.Error("error connecting:", err.Error())
				return
			}
			// Call connect callback
			b.onConnect(scanResult, device)
			//adapter.StopScan()
		}()
	}

}

func (b *Ble) onConnect(scanResult bluetooth.ScanResult, device bluetooth.Device) {
	println("connected:", scanResult.Address.String(), scanResult.LocalName())

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

	// buffer to retrieve characteristic data
	buf := make([]byte, 255)

	// Iterate services
	for _, service := range services {
		log.Infof("service: %v", service.String())

		// Get a list of characteristics below the service
		characteristics, err := service.DiscoverCharacteristics(nil)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		// Iterate characteristics
		for _, char := range characteristics {
			data := make([]byte, 10)
			n, _ := char.Read(data)
			log.Infof("char: %s, %v, %v", char.UUID().String(), n, data)
			mtu, err := char.GetMTU()
			if err != nil {
				println("    mtu: error:", err.Error())
			} else {
				println("    mtu:", mtu)
			}
			n, err = char.Read(buf)
			if err != nil {
				println("    ", err.Error())
			} else {
				println("    data bytes", strconv.Itoa(n))
				println("    value =", string(buf[:n]))
			}

			if err = char.EnableNotifications(b.handler); err != nil {
				log.Error(err.Error())
			}

			//if _, err = char.Write([]byte("on")); err != nil {
			//	log.Error(err.Error())
			//}
		}

	}
}

func (b *Ble) handler(buf []byte) {
	log.Infof("handler: %v, %s", buf, string(buf))
}
