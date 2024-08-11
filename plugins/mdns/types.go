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

package mdns

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/supervisor"
)

const (
	// Name ...
	Name = "mdns"
	// FuncEntityAction ...
	FuncEntityAction = "entityAction"

	Version = "0.0.1"

	AttrInstance = "instance"
	AttrService  = "service"
	AttrDomain   = "domain"
	AttrHost     = "host"
	AttrPort     = "port"
	AttrIpAddr   = "ipAddr"
	AttrText     = "text"
	ActionScan   = "SCAN"

	DefaultService = "_workstation._tcp"
	DefaultDomain  = "local"
)

// NewSettings ...
func NewSettings() map[string]*m.Attribute {
	return map[string]*m.Attribute{
		AttrInstance: {
			Name:  AttrInstance,
			Type:  common.AttributeString,
			Value: "escapepod",
		},
		AttrService: {
			Name:  AttrService,
			Type:  common.AttributeString,
			Value: DefaultService,
		},
		AttrDomain: {
			Name:  AttrDomain,
			Type:  common.AttributeString,
			Value: DefaultDomain,
		},
		AttrIpAddr: {
			Name:  AttrIpAddr,
			Type:  common.AttributeString,
			Value: "",
		},
		AttrText: {
			Name:  AttrText,
			Type:  common.AttributeString,
			Value: "txtv=0, lo=1, la=2",
		},
		AttrHost: {
			Name:  AttrHost,
			Type:  common.AttributeString,
			Value: "escapepod",
		},
		AttrPort: {
			Name:  AttrPort,
			Type:  common.AttributeInt,
			Value: 8084,
		},
	}
}

func NewActions() map[string]supervisor.ActorAction {
	return map[string]supervisor.ActorAction{
		ActionScan: {
			Name:        ActionScan,
			Description: "Scanning for devices begins. It stops after 10 seconds.",
		},
	}
}
