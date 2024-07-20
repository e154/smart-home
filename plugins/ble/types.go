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
	"github.com/e154/smart-home/system/supervisor"
)

const (
	// Name ...
	Name = "ble"
	// FuncEntityAction ...
	FuncEntityAction = "entityAction"

	Version = "0.0.1"

	AttrAddress              = "address"
	AttrTimeoutSec           = "timeout_sec"
	AttrConnectionTimeoutSec = "connection_timeout_sec"
	ActionScan               = "SCAN"
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
			Value: 5,
		},
		AttrConnectionTimeoutSec: {
			Name:  AttrConnectionTimeoutSec,
			Type:  common.AttributeInt,
			Value: 5,
		},
	}
}

func NewActions() map[string]supervisor.ActorAction {
	return map[string]supervisor.ActorAction{
		ActionScan: {
			Name:        ActionScan,
			Description: "scan",
		},
	}
}
