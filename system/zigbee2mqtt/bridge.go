// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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
	"context"
	"encoding/json"
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/metrics"
	"github.com/e154/smart-home/system/mqtt"
	"strings"
	"sync"
	"time"
)

const (
	homeassistantTopic = "homeassistant"
)

// Bridge ...
type Bridge struct {
	metric         *metrics.MetricManager
	adaptors       *adaptors.Adaptors
	mqtt           mqtt.MqttServ
	mqttClient     mqtt.MqttCli
	isStarted      bool
	settingsLock   sync.Mutex
	state          string // online|offline
	config         BridgeConfig
	devicesLock    sync.Mutex
	devices        map[string]*Device
	modelLock      sync.Mutex
	model          *m.Zigbee2mqtt
	networkmapLock sync.Mutex
	scanInProcess  bool
	lastScan       *time.Time
	networkmap     string
}

// NewBridge ...
func NewBridge(mqtt mqtt.MqttServ,
	adaptors *adaptors.Adaptors,
	model *m.Zigbee2mqtt,
	metric *metrics.MetricManager) *Bridge {
	return &Bridge{
		adaptors: adaptors,
		devices:  make(map[string]*Device),
		model:    model,
		metric:   metric,
		mqtt:     mqtt,
	}
}

// Start ...
func (g *Bridge) Start() {

	if g.isStarted {
		return
	}
	g.isStarted = true

	log.Infof("bridge id %v,  base topic: %v", g.model.Id, g.model.BaseTopic)

	g.metric.Update(metrics.Zigbee2MqttAdd{
		TotalNum: int64(len(g.model.Devices)),
	})

	if g.mqttClient == nil {
		g.mqttClient = g.mqtt.NewClient(fmt.Sprintf("bridge_%v", g.model.Name))
	}

	// /zigbee2mqtt/bridge/#
	g.mqttClient.Subscribe(fmt.Sprintf("%s/bridge/#", g.model.BaseTopic), g.onBridgePublish)

	// /homeassistant/#
	g.mqttClient.Subscribe(fmt.Sprintf("%s/#", homeassistantTopic), g.onAssistPublish)

	if err := g.safeGetDeviceList(); err != nil {
		log.Error(err.Error())

	}

	g.configPermitJoin(g.model.PermitJoin)
}

// Stop ...
func (g *Bridge) Stop(ctx context.Context) {
	if !g.isStarted {
		return
	}
	g.isStarted = false
	g.mqttClient.UnsubscribeAll()
}

func (g *Bridge) onBridgePublish(client mqtt.MqttCli, message mqtt.Message) {

	var topic = strings.Split(message.Topic, "/")

	switch topic[2] {
	case "state":
		g.onBridgeStatePublish(client, message)
	case "log", "logging":
		g.onLogPublish(client, message)
	case "config":
		g.onConfigPublish(client, message)
	case "networkmap":
		g.onNetworkmapPublish(client, message)
	default:
		log.Warnf("unknown topic %v", topic)
	}
}

func (g *Bridge) onAssistPublish(client mqtt.MqttCli, message mqtt.Message) {

	var topic = strings.Split(message.Topic, "/")

	// hemeassistant/sensor/0x00158d00031c8ef3/click/config
	// hemeassistant/sensor/0x00158d00031c8ef3/battery/config
	// hemeassistant/sensor/0x00158d00031c8ef3/linkquality/config

	deviceType := topic[1]
	friendlyName := topic[2]
	function := topic[3]
	deviceInfo := AssistDevice{}
	_ = json.Unmarshal(message.Payload, &deviceInfo)

	device, err := g.safeGetDevice(friendlyName)
	if err != nil {
		return
	}

	device.AddFunc(function)
	device.DeviceType(deviceType)
	device.SetStatus(active)

	if err = g.safeUpdateDevice(device); err != nil {
		log.Error(err.Error())
	}
}

func (g *Bridge) onBridgeStatePublish(client mqtt.MqttCli, message mqtt.Message) {
	g.settingsLock.Lock()
	g.state = string(message.Payload)
	g.settingsLock.Unlock()
}

func (g *Bridge) onNetworkmapPublish(client mqtt.MqttCli, message mqtt.Message) {

	var topic = strings.Split(message.Topic, "/")

	if len(topic) < 4 {
		return
	}

	switch topic[3] {
	case "raw":
		log.Info("method not implemented")
	case "graphviz":
		g.networkmapLock.Lock()
		g.scanInProcess = false
		g.lastScan = common.Time(time.Now())
		g.networkmap = string(message.Payload)
		g.networkmapLock.Unlock()
	}
}

func (g *Bridge) onConfigPublish(client mqtt.MqttCli, message mqtt.Message) {

	var topic = strings.Split(message.Topic, "/")

	if len(topic) > 3 {
		switch topic[3] {
		case "devices":
			g.onConfigDevicesPublish(client, message)
			return
		}
	}

	config := BridgeConfig{}
	_ = json.Unmarshal(message.Payload, &config)
	g.settingsLock.Lock()
	g.config = config
	g.settingsLock.Unlock()
}

func (g *Bridge) onConfigDevicesPublish(client mqtt.MqttCli, message mqtt.Message) {}

func (g *Bridge) safeGetDevice(friendlyName string) (device *Device, err error) {

	g.devicesLock.Lock()
	defer g.devicesLock.Unlock()

	//log.Infof("Get device %v", friendlyName)

	var ok bool
	if device, ok = g.devices[friendlyName]; !ok {
		err = common.ErrNotFound
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
		device.GetImage()

	} else {
		log.Infof("add device %v ...", model.Id)
		if err = g.adaptors.Zigbee2mqttDevice.Add(&model); err != nil {
			return
		}
		g.metric.Update(metrics.Zigbee2MqttAdd{TotalNum: 1})
	}

	g.devices[device.friendlyName] = device

	//TODO optimize
	g.modelLock.Lock()
	g.model.Devices = make([]*m.Zigbee2mqttDevice, len(g.devices))
	i := 0
	for _, device := range g.devices {
		g.model.Devices[i] = &device.model
		i++
	}
	g.modelLock.Unlock()

	g.metric.Update(metrics.Zigbee2MqttUpdate{})

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

func (g *Bridge) onLogPublish(client mqtt.MqttCli, message mqtt.Message) {
	var lm BridgeLog
	_ = json.Unmarshal(message.Payload, &lm)
	log.Infof("%s, %v, %s", lm.Message, lm.Meta, lm.Type)
	switch lm.Type {
	case "device_removed":
		g.deviceRemoved(lm.Message)
	case "device_force_removed":
		g.deviceForceRemoved(lm.Message)
	case "pairing":
		params := BridgePairingMeta{}
		_ = common.Copy(&params, lm.Meta, common.JsonEngine)
		g.devicePairing(params)
	}
}

func (g *Bridge) getState() string {
	g.settingsLock.Lock()
	defer g.settingsLock.Unlock()
	return g.state
}

// get config
func (g *Bridge) getConfig() {
	g.mqttClient.Publish(g.topic("/bridge/config/get"), []byte{})
}

// get device list
func (g *Bridge) getDevices() {
	g.mqttClient.Publish(g.topic("/bridge/config/devices/get"), []byte{})
}

func (g *Bridge) configPermitJoin(tr bool) {
	var permitJoin = "true"
	if !tr {
		permitJoin = "false"
	}
	g.mqttClient.Publish(g.topic("/bridge/config/permit_join"), []byte(permitJoin))
}

func (g *Bridge) configLastSeen() {}
func (g *Bridge) configElapsed()  {}

// Resets the ZNP (CC2530/CC2531).
func (g *Bridge) configReset() {
	g.mqttClient.Publish(g.topic("/bridge/config/reset"), []byte{})
}

// ConfigReset ...
func (g *Bridge) ConfigReset() {
	g.configReset()
	time.Sleep(time.Second * 5)
}

func (g *Bridge) configLogLevel() {}

// DeviceOptions ...
func (g *Bridge) DeviceOptions() {}

// Remove ...
func (g *Bridge) Remove(friendlyName string) {
	g.mqttClient.Publish(g.topic("/bridge/config/remove"), []byte(friendlyName))
}

// Ban ...
func (g *Bridge) Ban(friendlyName string) {
	g.mqttClient.Publish(g.topic("/bridge/config/force_remove"), []byte(friendlyName))
}

// Whitelist ...
func (g *Bridge) Whitelist(friendlyName string) {
	g.mqttClient.Publish(g.topic("/bridge/config/whitelist"), []byte(friendlyName))
}

// RenameDevice ...
func (g *Bridge) RenameDevice(friendlyName, name string) (err error) {

	var device *Device
	if device, err = g.safeGetDevice(friendlyName); err != nil {
		return
	}

	device.SetName(name)

	err = g.safeUpdateDevice(device)

	return
}

// RenameLast ...
func (g *Bridge) RenameLast() {}

// AddGroup ...
func (g *Bridge) AddGroup() {}

// RemoveGroup ...
func (g *Bridge) RemoveGroup() {}

// UpdateNetworkmap ...
func (g *Bridge) UpdateNetworkmap() {
	g.networkmapLock.Lock()
	defer g.networkmapLock.Unlock()

	if g.scanInProcess {
		return
	}
	g.scanInProcess = true

	g.mqttClient.Publish(g.topic("/bridge/networkmap"), []byte("graphviz"))
}

// Networkmap ...
func (g *Bridge) Networkmap() string {
	g.networkmapLock.Lock()
	defer g.networkmapLock.Unlock()
	return g.networkmap
}

func (g *Bridge) topic(s string) string {
	return fmt.Sprintf("%s%s", g.model.BaseTopic, s)
}

// GetDeviceTopic ...
func (g *Bridge) GetDeviceTopic(friendlyName string) string {
	return g.topic("/" + friendlyName)
}

func (g *Bridge) deviceRemoved(friendlyName string) {
	device, err := g.safeGetDevice(friendlyName)
	if err != nil {
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

func (g *Bridge) devicePairing(params BridgePairingMeta) {

	device, err := g.safeGetDevice(params.FriendlyName)
	if err != nil && err.Error() != "record not found" {
		log.Error(err.Error())
		return
	}

	if err != nil && err.Error() == "record not found" {
		model := &m.Zigbee2mqttDevice{
			Id:            params.FriendlyName,
			Status:        active,
			Zigbee2mqttId: g.model.Id,
			Name:          params.FriendlyName,
			Model:         params.Model,
			Description:   params.Description,
			Manufacturer:  params.Vendor,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		model.GetImageUrl()

		device = NewDevice(params.FriendlyName, model)
	}

	device.SetStatus(active)
	device.SetModel(params.Model)
	device.SetDescription(params.Description)
	device.SetVendor(params.Vendor)

	if err = g.safeUpdateDevice(device); err != nil {
		log.Error(err.Error())
	}
}

// PermitJoin ...
func (g *Bridge) PermitJoin(permitJoin bool) {
	g.modelLock.Lock()
	defer g.modelLock.Unlock()

	g.model.PermitJoin = permitJoin
	if err := g.adaptors.Zigbee2mqtt.Update(g.model); err != nil {
		return
	}
	g.configPermitJoin(g.model.PermitJoin)
}

// UpdateModel ...
func (g *Bridge) UpdateModel(model *m.Zigbee2mqtt) {
	g.modelLock.Lock()
	defer g.modelLock.Unlock()

	g.model.Login = model.Login
	g.model.BaseTopic = model.BaseTopic
	g.model.PermitJoin = model.PermitJoin
	g.model.EncryptedPassword = model.EncryptedPassword

	g.configPermitJoin(g.model.PermitJoin)
}

// Info ...
func (g *Bridge) Info() (info *Zigbee2mqttInfo) {

	g.networkmapLock.Lock()
	g.settingsLock.Lock()
	g.devicesLock.Lock()

	defer func() {
		g.networkmapLock.Unlock()
		g.settingsLock.Unlock()
		g.devicesLock.Unlock()
	}()

	model := m.Zigbee2mqtt{}

	_ = common.Copy(&model, g.model, common.JsonEngine)

	info = &Zigbee2mqttInfo{
		ScanInProcess: g.scanInProcess,
		LastScan:      g.lastScan,
		Networkmap:    g.networkmap,
		Status:        g.state,
		Model:         model,
	}

	return
}

// GetModel ...
func (g *Bridge) GetModel() (model m.Zigbee2mqtt) {
	g.modelLock.Lock()
	defer g.modelLock.Unlock()

	model = m.Zigbee2mqtt{}
	_ = common.Copy(&model, &g.model, common.JsonEngine)

	return
}
