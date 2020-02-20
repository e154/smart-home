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
	"strings"
	"sync"
)

const (
	baseTopic          = "zigbee2mqtt"
	homeassistantTopic = "homeassistant"
)

type Bridge struct {
	adaptors     *adaptors.Adaptors
	mqtt         *mqttServer.Mqtt
	mqttClient   *mqtt_client.Client
	isStarted    bool
	settingsLock sync.Mutex
	state        string // online|offline
	config       BridgeConfig
	devicesLock  sync.Mutex
	devices      map[string]Device
}

func NewBridge(mqtt *mqttServer.Mqtt,
	adaptors *adaptors.Adaptors) *Bridge {
	return &Bridge{
		mqtt:     mqtt,
		adaptors: adaptors,
		devices:  make(map[string]Device),
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
	if err := g.mqttClient.Subscribe(fmt.Sprintf("%s/bridge/#", baseTopic), 0, g.onBridgePublish); err != nil {
		log.Warning(err.Error())
	}

	// /homeassistant/#
	if err := g.mqttClient.Subscribe(fmt.Sprintf("%s/#", homeassistantTopic), 0, g.onAssistPublish); err != nil {
		log.Warning(err.Error())
	}

	if err := g.safeGetDeviceList(); err != nil {
		log.Error(err.Error())

	}

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
		model := &m.Zigbee2mqttDevice{
			Id:           friendlyName,
			Name:         friendlyName,
			Type:         deviceType,
			Functions:    []string{function},
			Model:        deviceInfo.Device.Model,
			Manufacturer: deviceInfo.Device.Manufacturer,
		}
		device = NewDevice(friendlyName, model)

	}

	if err == nil {
		device.AddFunc(function)
		device.DeviceType(deviceType)
		device.DeviceInfo(deviceInfo)
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
	config := BridgeConfig{}
	_ = json.Unmarshal(message.Payload(), &config)
	g.settingsLock.Lock()
	g.config = config
	g.settingsLock.Unlock()
}

func (g *Bridge) safeGetDevice(friendlyName string) (device Device, err error) {

	g.devicesLock.Lock()
	defer g.devicesLock.Unlock()

	var ok bool
	if device, ok = g.devices[friendlyName]; ok {
		return
	}

	log.Infof("Get device %v", friendlyName)

	var model *m.Zigbee2mqttDevice
	if model, err = g.adaptors.Zigbee2mqttDevice.GetById(friendlyName); err != nil {
		return
	}

	device = NewDevice(friendlyName, model)

	g.devices[friendlyName] = device

	return
}

func (g *Bridge) safeUpdateDevice(device Device) (err error) {
	g.devicesLock.Lock()
	defer g.devicesLock.Unlock()

	log.Infof("Update device %v", device.friendlyName)

	model := device.GetModel()

	if _, err = g.adaptors.Zigbee2mqttDevice.GetById(device.friendlyName); err == nil {
		if err = g.adaptors.Zigbee2mqttDevice.Update(&model); err != nil {
			log.Error(err.Error())
			return
		}
	} else {
		if err = g.adaptors.Zigbee2mqttDevice.Add(&model); err != nil {
			log.Error(err.Error())
			return
		}
	}

	g.devices[device.friendlyName] = device

	return
}

func (g *Bridge) safeGetDeviceList() (err error) {

	var models []*m.Zigbee2mqttDevice
	if models, _, err = g.adaptors.Zigbee2mqttDevice.List(999, 0); err != nil {
		return
	}

	g.devicesLock.Lock()
	defer g.devicesLock.Unlock()

	for _, model := range models {
		g.devices[model.Id] = NewDevice(model.Id, model)
	}

	return
}

func (g *Bridge) onLogPublish(client mqtt.Client, message mqtt.Message) {
	var lm BridgeLog
	_ = json.Unmarshal(message.Payload(), &lm)
	log.Infof("%v, %v, %v", lm.Type, lm.Message, lm.Meta)
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

func (g *Bridge) configPermitJoin() {}
func (g *Bridge) configLastSeen()   {}
func (g *Bridge) configElapsed()    {}

// Resets the ZNP (CC2530/CC2531).
func (g *Bridge) configReset() {
	g.mqtt.Publish(g.topic("/bridge/config/reset"), []byte{}, 0, false)
}

func (g *Bridge) configLogLevel() {}
func (g *Bridge) DeviceOptions()  {}
func (g *Bridge) Remove()         {}
func (g *Bridge) Ban()            {}
func (g *Bridge) Whitelist()      {}
func (g *Bridge) Rename()         {}
func (g *Bridge) RenameLast()     {}
func (g *Bridge) AddGroup()       {}
func (g *Bridge) RemoveGroup()    {}
func (g *Bridge) Networkmap()     {}

func (g *Bridge) topic(s string) string {
	return fmt.Sprintf("%s/%s", baseTopic, s)
}
