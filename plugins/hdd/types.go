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

package hdd

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/supervisor"
)

const (
	// EntityHDD ...
	EntityHDD = string("hdd")

	// AttrPath ...
	AttrPath = "path"

	// AttrFstype ...
	AttrFstype = "fstype"

	// AttrTotal ...
	AttrTotal = "total"

	// AttrFree ...
	AttrFree = "free"

	// AttrUsed ...
	AttrUsed = "used"

	// AttrUsedPercent ...
	AttrUsedPercent = "used_percent"

	// AttrInodesTotal ...
	AttrInodesTotal = "inodes_total"

	// AttrInodesUsed ...
	AttrInodesUsed = "inodes_used"

	// AttrInodesFree ...
	AttrInodesFree = "inodes_free"

	// AttrInodesUsedPercent ...
	AttrInodesUsedPercent = "inodes_sed_percent"

	// AttrMountPoint ...
	AttrMountPoint = "mount_point"

	// Name ...
	Name = "hdd"

	// EntityType ...
	EntityType = "hdd"

	ActionCheck = "CHECK"

	Version = "0.0.1"
)

// NewAttr ...
func NewAttr() m.Attributes {
	return m.Attributes{
		AttrPath: {
			Name: AttrPath,
			Type: common.AttributeString,
		},
		AttrFstype: {
			Name: AttrFstype,
			Type: common.AttributeString,
		},
		AttrTotal: {
			Name: AttrTotal,
			Type: common.AttributeInt,
		},
		AttrFree: {
			Name: AttrFree,
			Type: common.AttributeInt,
		},
		AttrUsed: {
			Name: AttrUsed,
			Type: common.AttributeInt,
		},
		AttrUsedPercent: {
			Name: AttrUsedPercent,
			Type: common.AttributeFloat,
		},
		AttrInodesTotal: {
			Name: AttrInodesTotal,
			Type: common.AttributeInt,
		},
		AttrInodesUsed: {
			Name: AttrInodesUsed,
			Type: common.AttributeInt,
		},
		AttrInodesFree: {
			Name: AttrInodesFree,
			Type: common.AttributeInt,
		},
		AttrInodesUsedPercent: {
			Name: AttrInodesUsedPercent,
			Type: common.AttributeFloat,
		},
	}
}

// NewSettings ...
func NewSettings() map[string]*m.Attribute {
	return map[string]*m.Attribute{
		AttrMountPoint: {
			Name: AttrMountPoint,
			Type: common.AttributeString,
		},
	}
}

// entity action list
func NewActions() map[string]supervisor.ActorAction {
	return map[string]supervisor.ActorAction{
		ActionCheck: {
			Name:        ActionCheck,
			Description: "check disk",
		},
	}
}
