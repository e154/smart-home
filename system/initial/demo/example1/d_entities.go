// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

package example1

import (
	"context"
	"fmt"
	"os"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/cgminer"
	"github.com/e154/smart-home/plugins/cgminer/bitmine"
	"github.com/e154/smart-home/plugins/sensor"
	"github.com/e154/smart-home/plugins/telegram"
	. "github.com/e154/smart-home/system/initial/assertions"
)

type Supervisor struct {
	adaptors *adaptors.Adaptors
}

func NewSupervisor(adaptors *adaptors.Adaptors) *Supervisor {
	return &Supervisor{
		adaptors: adaptors,
	}
}

func (e *Supervisor) addEntities(ctx context.Context, scripts []*m.Script, area *m.Area) (entities []*m.Entity, err error) {

	var script *m.Script
	if len(scripts) > 0 {
		script = scripts[0]
	}

	entity1 := e.addL3(ctx, "l3n1", "192.168.0.247", script, area)
	entity2 := e.addL3(ctx, "l3n2", "192.168.0.242", script, area)
	entity3 := e.addL3(ctx, "l3n3", "192.168.0.244", script, area)
	entity4 := e.addL3(ctx, "l3n4", "192.168.0.243", script, area)
	entity5 := e.addL3(ctx, "l3n5", "192.168.0.240", script, area)

	tgBot := e.addTgBot(ctx, "clavicus", os.Getenv("SH_TG_BOT_TOKEN"), script, area)
	sensorEntity := e.addSensor(ctx, "api", scripts[1], area)

	entities = []*m.Entity{entity1, entity2, entity3, entity4, entity5, tgBot, sensorEntity}

	return
}

func (e *Supervisor) addL3(ctx context.Context, name, host string, script *m.Script, area *m.Area) (ent *m.Entity) {

	id := common.EntityId(fmt.Sprintf("cgminer.%s", name))

	var err error
	if ent, err = e.adaptors.Entity.GetById(ctx, id); err == nil {
		return
	}

	settings := cgminer.NewSettings()
	settings[cgminer.SettingHost].Value = host
	settings[cgminer.SettingPort].Value = 4028
	settings[cgminer.SettingTimeout].Value = 2
	settings[cgminer.SettingUser].Value = "user"
	settings[cgminer.SettingPass].Value = "pass"
	settings[cgminer.SettingManufacturer].Value = bitmine.ManufactureBitmine
	settings[cgminer.SettingModel].Value = bitmine.DeviceL3Plus
	ent = &m.Entity{
		Id:          id,
		Description: "antminer L3+",
		PluginName:  cgminer.Name,
		AutoLoad:    true,
		Attributes:  cgminer.NewAttr(),
		Settings:    settings,
		Area:        area,
	}
	ent.Actions = []*m.EntityAction{
		{
			Name:        "CHECK",
			Description: "condition check",
			Script:      script,
		},
	}
	ent.States = []*m.EntityState{
		{
			Name:        "ENABLED",
			Description: "enabled state",
		},
		{
			Name:        "DISABLED",
			Description: "disabled state",
		},
		{
			Name:        "ERROR",
			Description: "error state",
		},
		{
			Name:        "WARNING",
			Description: "warning state",
		},
	}
	ent.Attributes = m.Attributes{
		"heat": {
			Name: "heat",
			Type: common.AttributeBool,
		},
		"chain1_temp_chip": {
			Name: "chain1_temp_chip",
			Type: common.AttributeInt,
		},
		"chain2_temp_chip": {
			Name: "chain2_temp_chip",
			Type: common.AttributeInt,
		},
		"chain3_temp_chip": {
			Name: "chain3_temp_chip",
			Type: common.AttributeInt,
		},
		"chain4_temp_chip": {
			Name: "chain4_temp_chip",
			Type: common.AttributeInt,
		},
		"chain1_temp_pcb": {
			Name: "chain1_temp_pcb",
			Type: common.AttributeInt,
		},
		"chain2_temp_pcb": {
			Name: "chain2_temp_pcb",
			Type: common.AttributeInt,
		},
		"chain3_temp_pcb": {
			Name: "chain3_temp_pcb",
			Type: common.AttributeInt,
		},
		"chain4_temp_pcb": {
			Name: "chain4_temp_pcb",
			Type: common.AttributeInt,
		},
		"chain_acn1": {
			Name: "chain_acn1",
			Type: common.AttributeInt,
		},
		"chain_acn2": {
			Name: "chain_acn2",
			Type: common.AttributeInt,
		},
		"chain_acn3": {
			Name: "chain_acn3",
			Type: common.AttributeInt,
		},
		"chain_acn4": {
			Name: "chain_acn4",
			Type: common.AttributeInt,
		},
		"fan1": {
			Name: "fan1",
			Type: common.AttributeInt,
		},
		"fan2": {
			Name: "fan2",
			Type: common.AttributeInt,
		},
		"ghs_av": {
			Name: "ghs_av",
			Type: common.AttributeInt,
		},
		"hardware_errors": {
			Name: "hardware_errors",
			Type: common.AttributeInt,
		},
	}

	err = e.adaptors.Entity.Add(ctx, ent)
	So(err, ShouldBeNil)

	_, err = e.adaptors.EntityStorage.Add(ctx, &m.EntityStorage{
		EntityId:   ent.Id,
		Attributes: ent.Attributes.Serialize(),
	})
	So(err, ShouldBeNil)

	return
}

func (e *Supervisor) addTgBot(ctx context.Context, name, token string, script *m.Script, area *m.Area) (ent *m.Entity) {

	id := common.EntityId(fmt.Sprintf("%s.%s", telegram.Name, name))

	var err error
	if ent, err = e.adaptors.Entity.GetById(ctx, id); err == nil {
		return
	}

	settings := telegram.NewSettings()
	settings[telegram.AttrToken].Value = token
	ent = &m.Entity{
		Id:          id,
		Description: "",
		PluginName:  telegram.Name,
		AutoLoad:    true,
		Attributes:  telegram.NewAttr(),
		Settings:    settings,
		Actions: []*m.EntityAction{
			{
				Name:        "CHECK",
				Description: "check status",
				Script:      script,
			},
		},
		Area: area,
	}
	err = e.adaptors.Entity.Add(ctx, ent)
	So(err, ShouldBeNil)
	_, err = e.adaptors.EntityStorage.Add(ctx, &m.EntityStorage{
		EntityId:   ent.Id,
		Attributes: ent.Attributes.Serialize(),
	})
	So(err, ShouldBeNil)

	return
}

func (e *Supervisor) addSensor(ctx context.Context, name string, script *m.Script, area *m.Area) (ent *m.Entity) {

	id := common.EntityId(fmt.Sprintf("cgminer.%s", name))

	var err error
	if ent, err = e.adaptors.Entity.GetById(ctx, id); err == nil {
		return
	}

	ent = &m.Entity{
		Id:          id,
		Description: "",
		PluginName:  sensor.Name,
		AutoLoad:    true,
		Attributes: m.Attributes{
			"paid_rewards": {
				Name: "paid_rewards",
				Type: common.AttributeFloat,
			},
		},
		Actions: []*m.EntityAction{
			{
				Name:        "CHECK",
				Description: "condition check",
				Script:      script,
			},
		},
		States: []*m.EntityState{
			{
				Name:        "ENABLED",
				Description: "enabled state",
			},
			{
				Name:        "DISABLED",
				Description: "disabled state",
			},
			{
				Name:        "ERROR",
				Description: "error state",
			},
		},
		Area: area,
	}
	err = e.adaptors.Entity.Add(ctx, ent)
	So(err, ShouldBeNil)
	_, err = e.adaptors.EntityStorage.Add(ctx, &m.EntityStorage{
		EntityId:   ent.Id,
		Attributes: ent.Attributes.Serialize(),
	})
	So(err, ShouldBeNil)

	return
}
