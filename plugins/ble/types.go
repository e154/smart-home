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

package ble

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/triggers"
	"github.com/e154/smart-home/system/supervisor"
)

const (
	// Name ...
	Name = "ble"
	// FuncEntityAction ...
	FuncEntityAction = "entityAction"

	Version = "0.0.1"

	AttrAddress              = "address"
	AttrCharacteristic       = "characteristic"
	AttrService              = "service"
	AttrTimeoutSec           = "timeout_sec"
	AttrConnectionTimeoutSec = "connection_timeout_sec"
	AttrDebug                = "debug"
	ActionScan               = "SCAN"
	AttrSystemInfo           = "SystemInfo"

	FunctionName = "automationTriggerBle"

	DefaultTimeout           int64 = 5
	DefaultConnectionTimeout int64 = 5
)

// NewSettings ...
func NewSettings() map[string]*m.Attribute {
	return map[string]*m.Attribute{
		AttrAddress: {
			Name:  AttrAddress,
			Type:  common.AttributeString,
			Value: "",
		},
		AttrTimeoutSec: {
			Name:  AttrTimeoutSec,
			Type:  common.AttributeInt,
			Value: DefaultTimeout,
		},
		AttrConnectionTimeoutSec: {
			Name:  AttrConnectionTimeoutSec,
			Type:  common.AttributeInt,
			Value: DefaultConnectionTimeout,
		},
		AttrDebug: {
			Name:  AttrDebug,
			Type:  common.AttributeBool,
			Value: false,
		},
	}
}

func NewActions() map[string]supervisor.ActorAction {
	return map[string]supervisor.ActorAction{
		ActionScan: {
			Name:        ActionScan,
			Description: "Scan starts a BLE scan. It is stopped after 10 seconds.",
		},
	}
}

func NewTriggerParams() m.TriggerParams {
	return m.TriggerParams{
		Script:   true,
		Required: []string{AttrAddress, AttrCharacteristic},
		Attributes: m.Attributes{
			AttrSystemInfo: {
				Name: AttrSystemInfo,
				Type: common.AttributeNotice,
			},
			AttrAddress: {
				Name: AttrAddress,
				Type: common.AttributeString,
			},
			//AttrService: {
			//	Name: AttrService,
			//	Type: common.AttributeString,
			//},
			AttrCharacteristic: {
				Name: AttrCharacteristic,
				Type: common.AttributeString,
			},
		},
	}
}

type TriggerParams struct {
	Bluetooth
	options triggers.Subscriber
}

type Bluetooth interface {
	Connect() error
	IsConnected() bool
	Disconnect() error
	Scan(address *string)
	Write(char string, request []byte, withResponse bool) ([]byte, error)
	Read(char string) ([]byte, error)
	Subscribe(char string, handler func([]byte)) error
}
