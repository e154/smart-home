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

package zigbee2mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/e154/smart-home/internal/common"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/apperr"
	pkgCommon "github.com/e154/smart-home/pkg/common"
	"github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/mqtt"
)

// Bridge ...
type Bridge struct {
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
	model          *models.Zigbee2mqtt
	networkmapLock sync.Mutex
	scanInProcess  bool
	lastScan       *time.Time
	networkmap     string
}

// NewBridge ...
func NewBridge(mqtt mqtt.MqttServ,
	adaptors *adaptors.Adaptors,
	model *models.Zigbee2mqtt) *Bridge {
	return &Bridge{
		adaptors: adaptors,
		devices:  make(map[string]*Device),
		model:    model,
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

	//todo add metric ...

	if g.mqttClient == nil {
		g.mqttClient = g.mqtt.NewClient(fmt.Sprintf("bridge_%v", g.model.Name))
	}

	// /zigbee2mqtt/bridge/#
	_ = g.mqttClient.Subscribe(fmt.Sprintf("%s/bridge/#", g.model.BaseTopic), g.onBridgePublish)

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
	g.mqtt.RemoveClient(fmt.Sprintf("bridge_%v", g.model.Name))
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
	case "devices":
		g.onDevices(client, message)
	case "event":
		g.onEvent(client, message)
	default:
		log.Warnf("unknown topic %v", message.Topic)
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
		g.lastScan = pkgCommon.Time(time.Now())
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
		err = apperr.ErrNotFound
		return
	}

	return
}

func (g *Bridge) safeUpdateDevice(device *Device) (err error) {
	g.devicesLock.Lock()
	defer g.devicesLock.Unlock()

	model := device.GetModel()

	if _, err = g.adaptors.Zigbee2mqttDevice.GetById(context.Background(), device.friendlyName); err == nil {
		log.Infof("update device %v ...", model.Id)
		if err = g.adaptors.Zigbee2mqttDevice.Update(context.Background(), &model); err != nil {
			log.Error(err.Error())
			return
		}
		device.GetImage()

	} else {
		log.Infof("add device %v ...", model.Id)
		if err = g.adaptors.Zigbee2mqttDevice.Add(context.Background(), &model); err != nil {
			return
		}
		//todo add metric ...
	}

	g.devices[device.friendlyName] = device

	//TODO optimize
	g.modelLock.Lock()
	g.model.Devices = make([]*models.Zigbee2mqttDevice, len(g.devices))
	i := 0
	for _, device := range g.devices {
		g.model.Devices[i] = &device.model
		i++
	}
	g.modelLock.Unlock()

	//todo add metric ...

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
}

func (g *Bridge) onEvent(client mqtt.MqttCli, message mqtt.Message) {
	event := Event{}
	_ = json.Unmarshal(message.Payload, &event)
	switch event.Type {
	case EventDeviceAnnounce:
	case EventDeviceLeave:
		g.deviceLeave(event.Data.FriendlyName)
	case EventDeviceJoined:
		g.deviceJoined(event.Data.FriendlyName)
	case EventDeviceInterview:
		g.deviceInterview(event)
	default:
		log.Warnf("unknown event type \"%s\"", event.Type)

	}
}

func (g *Bridge) getState() string {
	g.settingsLock.Lock()
	defer g.settingsLock.Unlock()
	return g.state
}

// get config
func (g *Bridge) getConfig() {
	_ = g.mqttClient.Publish(g.topic("/bridge/config/get"), []byte{})
}

// get device list
func (g *Bridge) getDevices() {
	_ = g.mqttClient.Publish(g.topic("/bridge/config/devices/get"), []byte{})
}

func (g *Bridge) configPermitJoin(tr bool) {
	var permitJoin = "true"
	if !tr {
		permitJoin = "false"
	}
	_ = g.mqttClient.Publish(g.topic("/bridge/config/permit_join"), []byte(permitJoin))
}

func (g *Bridge) configLastSeen() {}
func (g *Bridge) configElapsed()  {}

// Resets the ZNP (CC2530/CC2531).
func (g *Bridge) configReset() {
	_ = g.mqttClient.Publish(g.topic("/bridge/config/reset"), []byte{})
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
	_ = g.mqttClient.Publish(g.topic("/bridge/config/remove"), []byte(friendlyName))
}

// Ban ...
func (g *Bridge) Ban(friendlyName string) {
	_ = g.mqttClient.Publish(g.topic("/bridge/config/force_remove"), []byte(friendlyName))
}

// Whitelist ...
func (g *Bridge) Whitelist(friendlyName string) {
	_ = g.mqttClient.Publish(g.topic("/bridge/config/whitelist"), []byte(friendlyName))
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

	_ = g.mqttClient.Publish(g.topic("/bridge/networkmap"), []byte("graphviz"))
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

func (g *Bridge) deviceLeave(friendlyName string) {
	device, err := g.safeGetDevice(friendlyName)
	if err != nil {
		return
	}
	device.SetStatus(removed)
	if err = g.safeUpdateDevice(device); err != nil {
		log.Error(err.Error())
	}
}

func (g *Bridge) deviceJoined(friendlyName string) {
	device, err := g.safeGetDevice(friendlyName)
	if err != nil {
		return
	}
	device.SetStatus(active)
	if err = g.safeUpdateDevice(device); err != nil {
		log.Error(err.Error())
	}
}

func (g *Bridge) deviceInterview(event Event) {

	log.Infof("device interview %s, status: %s", event.Data.FriendlyName, event.Data.Status)

	if event.Data.Status != "successful" {
		return
	}

	g.updateDevice(event.Data)
}

func (g *Bridge) onDevices(client mqtt.MqttCli, message mqtt.Message) {

	devices := make([]DeviceInfo, 0)
	_ = json.Unmarshal(message.Payload, &devices)

	for _, device := range devices {
		if device.Type == Coordinator || !device.InterviewCompleted {
			continue
		}
		g.updateDevice(device)

		go g.UpdateNetworkmap()
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

func (g *Bridge) updateDevice(params DeviceInfo) {

	payload, _ := json.Marshal(params)
	model := &models.Zigbee2mqttDevice{
		Id:            params.FriendlyName,
		Zigbee2mqttId: g.model.Id,
		Name:          params.FriendlyName,
		Type:          params.Type,
		Model:         params.Definition.Model,
		Description:   params.Definition.Description,
		Manufacturer:  params.Definition.Vendor,
		Status:        active,
		Payload:       payload,
	}

	device := NewDevice(params.FriendlyName, model)
	if err := g.safeUpdateDevice(device); err != nil {
		log.Error(err.Error())
	}
}

// PermitJoin ...
func (g *Bridge) PermitJoin(permitJoin bool) {
	g.modelLock.Lock()
	defer g.modelLock.Unlock()

	g.model.PermitJoin = permitJoin
	if err := g.adaptors.Zigbee2mqtt.Update(context.Background(), g.model); err != nil {
		return
	}
	g.configPermitJoin(g.model.PermitJoin)
}

// UpdateModel ...
func (g *Bridge) UpdateModel(model *models.Zigbee2mqtt) {
	g.modelLock.Lock()
	defer g.modelLock.Unlock()

	g.model.Name = model.Name
	g.model.BaseTopic = model.BaseTopic
	g.model.Login = model.Login
	g.model.EncryptedPassword = model.EncryptedPassword
	g.model.PermitJoin = model.PermitJoin

	g.configPermitJoin(g.model.PermitJoin)
}

// Info ...
func (g *Bridge) Info() (info *models.Zigbee2mqttInfo) {

	g.networkmapLock.Lock()
	g.settingsLock.Lock()
	g.devicesLock.Lock()

	defer func() {
		g.networkmapLock.Unlock()
		g.settingsLock.Unlock()
		g.devicesLock.Unlock()
	}()

	model := models.Zigbee2mqtt{}

	_ = common.Copy(&model, g.model, common.JsonEngine)

	info = &models.Zigbee2mqttInfo{
		ScanInProcess: g.scanInProcess,
		LastScan:      g.lastScan,
		Networkmap:    g.networkmap,
		Status:        g.state,
	}

	return
}

// GetModel ...
func (g *Bridge) GetModel() (model models.Zigbee2mqtt) {
	g.modelLock.Lock()
	defer g.modelLock.Unlock()

	model = models.Zigbee2mqtt{}
	_ = common.Copy(&model, &g.model, common.JsonEngine)

	return
}
