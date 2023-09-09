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
	"time"

	m "github.com/e154/smart-home/models"
)

// Zigbee2mqtt ...
type Zigbee2mqtt interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
	AddBridge(model *m.Zigbee2mqtt) (err error)
	GetBridgeById(id int64) (*m.Zigbee2mqtt, error)
	GetBridgeInfo(id int64) (*Zigbee2mqttBridge, error)
	ListBridges(limit, offset int64, order, sortBy string) (models []*Zigbee2mqttBridge, total int64, err error)
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

// Zigbee2mqttBridge ...
type Zigbee2mqttBridge struct {
	m.Zigbee2mqtt
	ScanInProcess bool       `json:"scan_in_process"`
	LastScan      *time.Time `json:"last_scan"`
	Networkmap    string     `json:"networkmap"`
	Status        string     `json:"status"`
}

// BridgeLog ...
type BridgeLog struct {
	Type    string                 `json:"type"`
	Message string                 `json:"message"`
	Meta    map[string]interface{} `json:"meta"`
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

const (
	active  = "active"
	banned  = "banned"
	removed = "removed"
)

const (
	// EventDeviceAnnounce ...
	EventDeviceAnnounce = "device_announce"
	// EventDeviceLeave ...
	EventDeviceLeave = "device_leave"
	// EventDeviceJoined ...
	EventDeviceJoined = "device_joined"
	// EventDeviceInterview ...
	EventDeviceInterview = "device_interview"
)

const (
	// StatusStarted ...
	StatusStarted = "started"
	// StatusFailed ...
	StatusFailed = "failed"
)

// EventDeviceInfoDefExpose ...
type EventDeviceInfoDefExpose struct {
	Access      int64                      `json:"access,omitempty"`
	Description string                     `json:"description,omitempty"`
	Name        string                     `json:"name,omitempty"`
	Property    string                     `json:"property,omitempty"`
	Type        string                     `json:"type,omitempty"`
	Unit        string                     `json:"unit,omitempty"`
	ValueMax    int64                      `json:"value_max,omitempty"`
	ValueMin    int64                      `json:"value_min,omitempty"`
	Values      []string                   `json:"values,omitempty"`
	Features    []EventDeviceInfoDefExpose `json:"features,omitempty"`
}

// EventDeviceInfoDef ...
type EventDeviceInfoDef struct {
	Description string                     `json:"description,omitempty"`
	Exposes     []EventDeviceInfoDefExpose `json:"exposes"`
	Model       string                     `json:"model"`
	Options     []string                   `json:"options"`
	SupportsOta bool                       `json:"supports_ota"`
	Vendor      string                     `json:"vendor"`
}

const (
	// Coordinator ...
	Coordinator = "Coordinator"
	// EndDevice ...
	EndDevice = "EndDevice"
)

// DeviceInfo ...
type DeviceInfo struct {
	Definition         EventDeviceInfoDef `json:"definition"`
	FriendlyName       string             `json:"friendly_name"`
	IeeeAddress        string             `json:"ieee_address"`
	Status             string             `json:"status,omitempty"`
	InterviewCompleted bool               `json:"interview_completed,omitempty"`
	Interviewing       bool               `json:"interviewing,omitempty"`
	ModelId            string             `json:"model_id,omitempty"`
	NetworkAddress     int64              `json:"network_address,omitempty"`
	PowerSource        string             `json:"power_source,omitempty"`
	Supported          bool               `json:"supported,omitempty"`
	Type               string             `json:"type,omitempty"`
}

// Event ...
type Event struct {
	Data DeviceInfo `json:"data"`
	Type string     `json:"type"`
}
