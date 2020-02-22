// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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

package zigbee2mqtt

import (
	"encoding/json"
	"fmt"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	mqttServer "github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/mqtt_client"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"regexp"
	"strings"
	"sync"
)

const (
	homeassistantTopic = "homeassistant"
)

type Bridge struct {
	model        *m.Zigbee2mqtt
	adaptors     *adaptors.Adaptors
	mqtt         *mqttServer.Mqtt
	mqttClient   *mqtt_client.Client
	isStarted    bool
	settingsLock sync.Mutex
	state        string // online|offline
	config       BridgeConfig
	devicesLock  sync.Mutex
	devices      map[string]*Device
}

func NewBridge(mqtt *mqttServer.Mqtt,
	adaptors *adaptors.Adaptors,
	model *m.Zigbee2mqtt) *Bridge {
	return &Bridge{
		mqtt:     mqtt,
		adaptors: adaptors,
		devices:  make(map[string]*Device),
		model:    model,
	}
}

func (g *Bridge) Start() {

	if g.isStarted {
		return
	}
	g.isStarted = true

	var err error
	if g.mqttClient == nil {
		log.Info("create new mqtt client...")
		if g.mqttClient, err = g.mqtt.NewClient(nil); err != nil {
			log.Error(err.Error())
		}
	}

	if !g.mqttClient.IsConnected() {
		if err = g.mqttClient.Connect(); err != nil {
			log.Error(err.Error())
		}
	}

	// /zigbee2mqtt/bridge/#
	if err := g.mqttClient.Subscribe(fmt.Sprintf("%s/bridge/#", g.model.BaseTopic), 0, g.onBridgePublish); err != nil {
		log.Warning(err.Error())
	}

	// /homeassistant/#
	if err := g.mqttClient.Subscribe(fmt.Sprintf("%s/#", homeassistantTopic), 0, g.onAssistPublish); err != nil {
		log.Warning(err.Error())
	}

	if err := g.safeGetDeviceList(); err != nil {
		log.Error(err.Error())

	}

	g.configPermitJoin(g.model.PermitJoin)

	log.Info("Starting...")
}

func (g *Bridge) Stop() {
	if !g.isStarted {
		return
	}
	g.isStarted = false
	g.mqttClient.UnsubscribeAll()
	g.mqttClient.Disconnect()
}

func (g *Bridge) onBridgePublish(client mqtt.Client, message mqtt.Message) {

	var topic = strings.Split(message.Topic(), "/")

	switch topic[2] {
	case "state":
		g.onBridgeStatePublish(client, message)
	case "log":
		g.onLogPublish(client, message)
	case "config":
		g.onConfigPublish(client, message)
	default:
		log.Warningf("unknown topic %v", topic)
	}
}

func (g *Bridge) onAssistPublish(client mqtt.Client, message mqtt.Message) {

	var topic = strings.Split(message.Topic(), "/")

	// hemeassistant/sensor/0x00158d00031c8ef3/click/config
	// hemeassistant/sensor/0x00158d00031c8ef3/battery/config
	// hemeassistant/sensor/0x00158d00031c8ef3/linkquality/config

	deviceType := topic[1]
	friendlyName := topic[2]
	function := topic[3]
	deviceInfo := AssistDevice{}
	_ = json.Unmarshal(message.Payload(), &deviceInfo)

	device, err := g.safeGetDevice(friendlyName)
	if err != nil && err.Error() != "record not found" {
		log.Error(err.Error())
		return
	}

	if err != nil && err.Error() == "record not found" {
		r := regexp.MustCompile(`\((.*?)\)`)
		devModel := r.FindString(deviceInfo.Device.Model)
		modelDesc := strings.Replace(deviceInfo.Device.Model, devModel, "", -1)
		devModel = strings.Replace(devModel, "(", "", -1)
		devModel = strings.Replace(devModel, ")", "", -1)

		model := &m.Zigbee2mqttDevice{
			Id:            friendlyName,
			Status:        active,
			Zigbee2mqttId: g.model.Id,
			Name:          friendlyName,
			Type:          deviceType,
			Functions:     []string{function},
			Model:         devModel,
			Description:   modelDesc,
			Manufacturer:  deviceInfo.Device.Manufacturer,
		}
		device = NewDevice(friendlyName, model)

	}

	if err == nil {
		device.AddFunc(function)
		device.DeviceType(deviceType)
		device.SetStatus(active)
	}

	if err = g.safeUpdateDevice(device); err != nil {
		log.Error(err.Error())
	}
}

func (g *Bridge) onBridgeStatePublish(client mqtt.Client, message mqtt.Message) {
	g.settingsLock.Lock()
	g.state = string(message.Payload())
	g.settingsLock.Unlock()
}

func (g *Bridge) onConfigPublish(client mqtt.Client, message mqtt.Message) {

	var topic = strings.Split(message.Topic(), "/")

	if len(topic) > 3 {
		switch topic[3] {
		case "devices":
			g.onConfigDevicesPublish(client, message)
			return
		}
	}

	config := BridgeConfig{}
	_ = json.Unmarshal(message.Payload(), &config)
	g.settingsLock.Lock()
	g.config = config
	g.settingsLock.Unlock()
}

func (g *Bridge) onConfigDevicesPublish(client mqtt.Client, message mqtt.Message) {

}

func (g *Bridge) safeGetDevice(friendlyName string) (device *Device, err error) {

	g.devicesLock.Lock()
	defer g.devicesLock.Unlock()

	log.Infof("Get device %v", friendlyName)

	var ok bool
	if device, ok = g.devices[friendlyName]; !ok {
		err = adaptors.ErrRecordNotFound
		return
	}

	return
}

func (g *Bridge) safeUpdateDevice(device *Device) (err error) {
	g.devicesLock.Lock()
	defer g.devicesLock.Unlock()

	model := device.GetModel()

	if _, err = g.adaptors.Zigbee2mqttDevice.GetById(device.friendlyName); err == nil {
		log.Infof("update device %v ...", model.Id)
		if err = g.adaptors.Zigbee2mqttDevice.Update(&model); err != nil {
			log.Error(err.Error())
			return
		}
	} else {
		log.Infof("add device %v ...", model.Id)
		if err = g.adaptors.Zigbee2mqttDevice.Add(&model); err != nil {
			log.Error(err.Error())
			return
		}
	}

	g.devices[device.friendlyName] = device

	return
}

func (g *Bridge) safeGetDeviceList() (err error) {

	g.devicesLock.Lock()
	defer g.devicesLock.Unlock()

	for _, model := range g.model.Devices {
		log.Infof("add device %v ...", model.Id)
		g.devices[model.Id] = NewDevice(model.Id, model)
	}

	return
}

func (g *Bridge) onLogPublish(client mqtt.Client, message mqtt.Message) {
	var lm BridgeLog
	_ = json.Unmarshal(message.Payload(), &lm)
	log.Infof("%v, %v, %v", lm.Type, lm.Message, lm.Meta)
	switch lm.Type {
	case "device_removed":
		g.deviceRemoved(lm.Message)
	case "device_force_removed":
		g.deviceForceRemoved(lm.Message)
	case "pairing":
		if friendlyName, ok := lm.Meta["friendly_name"].(string); ok {
			g.devicePairing(friendlyName)
		}
	}
}

func (g *Bridge) getState() string {
	g.settingsLock.Lock()
	defer g.settingsLock.Unlock()
	return g.state
}

// get config
func (g *Bridge) getConfig() {
	g.mqtt.Publish(g.topic("/bridge/config/get"), []byte{}, 0, false)
}

// get device list
func (g *Bridge) getDevices() {
	g.mqtt.Publish(g.topic("/bridge/config/devices/get"), []byte{}, 0, false)
}

func (g *Bridge) configPermitJoin(tr bool) {
	var permitJoin = "true"
	if !tr {
		permitJoin = "false"
	}
	g.mqtt.Publish(g.topic("/bridge/config/permit_join"), []byte(permitJoin), 0, false)
}

func (g *Bridge) configLastSeen() {}
func (g *Bridge) configElapsed()  {}

// Resets the ZNP (CC2530/CC2531).
func (g *Bridge) configReset() {
	g.mqtt.Publish(g.topic("/bridge/config/reset"), []byte{}, 0, false)
}

func (g *Bridge) configLogLevel() {}
func (g *Bridge) DeviceOptions()  {}

func (g *Bridge) Remove(friendlyName string) {
	g.mqtt.Publish(g.topic("/bridge/config/remove"), []byte(friendlyName), 0, false)
}

func (g *Bridge) Ban(friendlyName string) {
	g.mqtt.Publish(g.topic("/bridge/config/force_remove"), []byte(friendlyName), 0, false)
}

func (g *Bridge) Whitelist(friendlyName string) {
	g.mqtt.Publish(g.topic("/bridge/config/whitelist"), []byte(friendlyName), 0, false)
}

func (g *Bridge) Rename()      {}
func (g *Bridge) RenameLast()  {}
func (g *Bridge) AddGroup()    {}
func (g *Bridge) RemoveGroup() {}
func (g *Bridge) Networkmap()  {}

func (g *Bridge) topic(s string) string {
	return fmt.Sprintf("%s/%s", g.model.BaseTopic, s)
}

func (g *Bridge) GetDeviceTopic(friendlyName string) string {
	return g.topic(friendlyName)
}

func (g *Bridge) deviceRemoved(friendlyName string) {
	device, err := g.safeGetDevice(friendlyName)
	if err != nil {
		log.Error(err.Error())
		return
	}
	device.SetStatus(removed)
	if err = g.safeUpdateDevice(device); err != nil {
		log.Error(err.Error())
	}
}

func (g *Bridge) deviceForceRemoved(friendlyName string) {
	device, err := g.safeGetDevice(friendlyName)
	if err != nil {
		log.Error(err.Error())
		return
	}
	device.SetStatus(banned)
	if err = g.safeUpdateDevice(device); err != nil {
		log.Error(err.Error())
	}
}

func (g *Bridge) devicePairing(friendlyName string) {
	device, err := g.safeGetDevice(friendlyName)
	if err != nil {
		log.Error(err.Error())
		return
	}
	if device.Status() == active {
		return
	}
	device.SetStatus(active)
	if err = g.safeUpdateDevice(device); err != nil {
		log.Error(err.Error())
	}
}

func (g *Bridge) PermitJoin(permitJoin bool) {
	g.model.PermitJoin = permitJoin
	if err := g.adaptors.Zigbee2mqtt.Update(g.model); err != nil {
		return
	}
	g.configPermitJoin(g.model.PermitJoin)
}
