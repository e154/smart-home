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
	m "github.com/e154/smart-home/models"
	"time"
)

// Zigbee2mqtt ...
type Zigbee2mqtt interface {
	Start()
	Shutdown()
	AddBridge(model *m.Zigbee2mqtt) (err error)
	GetBridgeById(id int64) (*m.Zigbee2mqtt, error)
	GetBridgeInfo(id int64) (*Zigbee2mqttInfo, error)
	ListBridges(limit, offset int64, order, sortBy string) (models []*Zigbee2mqttInfo, total int64, err error)
	UpdateBridge(model *m.Zigbee2mqtt) (result *m.Zigbee2mqtt, err error)
	DeleteBridge(bridgeId int64) (err error)
	ResetBridge(bridgeId int64) (err error)
	BridgeDeviceBan(bridgeId int64, friendlyName string) (err error)
	BridgeDeviceWhitelist(bridgeId int64, friendlyName string) (err error)
	BridgeNetworkmap(bridgeId int64) (networkmap string, err error)
	BridgeUpdateNetworkmap(bridgeId int64) (err error)
	GetTopicByDevice(model *m.Zigbee2mqttDevice) (topic string, err error)
	DeviceRename(friendlyName, name string) (err error)
}

// DeviceType ...
type DeviceType string

// BridgeLog ...
type BridgeLog struct {
	Type    string                 `json:"type"`
	Message string                 `json:"message"`
	Meta    map[string]interface{} `json:"meta"`
}

// BridgePairingMeta ...
type BridgePairingMeta struct {
	FriendlyName string `json:"friendly_name"`
	Model        string `json:"model"`
	Vendor       string `json:"vendor"`
	Description  string `json:"description"`
	Supported    bool   `json:"supported"`
}

// BridgeConfigMeta ...
type BridgeConfigMeta struct {
	Transportrev int64 `json:"transportrev"`
	Product      int64 `json:"product"`
	Majorrel     int64 `json:"majorrel"`
	Minorrel     int64 `json:"minorrel"`
	Maintrel     int64 `json:"maintrel"`
	Revision     int64 `json:"revision"`
}

// BridgeConfigCoordinator ...
type BridgeConfigCoordinator struct {
	Type string           `json:"type"`
	Meta BridgeConfigMeta `json:"meta"`
}

// BridgeConfig ...
type BridgeConfig struct {
	Version     string                  `json:"version"`
	Commit      string                  `json:"commit"`
	Coordinator BridgeConfigCoordinator `json:"coordinator"`
	LogLevel    string                  `json:"log_level"`
	PermitJoin  string                  `json:"permit_join"`
}

// AssistDeviceInfo ...
type AssistDeviceInfo struct {
	Name         string `json:"name"`
	Model        string `json:"model"`
	Manufacturer string `json:"manufacturer"`
}

// AssistDevice ...
type AssistDevice struct {
	Device AssistDeviceInfo `json:"device"`
}

const (
	active  = "active"
	banned  = "banned"
	removed = "removed"
)

// Zigbee2mqttInfo ...
type Zigbee2mqttInfo struct {
	ScanInProcess bool          `json:"scan_in_process"`
	LastScan      *time.Time    `json:"last_scan"`
	Networkmap    string        `json:"networkmap"`
	Status        string        `json:"status"`
	Model         m.Zigbee2mqtt `json:"model"`
}
