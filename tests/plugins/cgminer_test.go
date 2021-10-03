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

package plugins

import (
	context2 "context"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/cgminer"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scripts"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestCgMiner(t *testing.T) {

	const cgminerSourceScript = `

ifError =(res)->
    return !res || res.error || res.Error

checkStatus =->
    stats = Miner.stats()
    if ifError(stats)
        Actor.setState
            'new_state': 'ERROR'
        return
    p = JSON.parse(stats.result)
    attrs = {
        heat: false
        chain1_temp_chip: p.temp2_1
        chain2_temp_chip: p.temp2_2
        chain3_temp_chip: p.temp2_3
        chain4_temp_chip: p.temp2_4
        chain1_temp_pcb: p.temp1
        chain2_temp_pcb: p.temp2
        chain3_temp_pcb: p.temp3
        chain4_temp_pcb: p.temp4
        fan1: p.fan1
        fan2: p.fan2
    }

    if p.temp1 > 60 || p.temp2 > 60 || p.temp3 > 60 || p.temp4 > 60 || p.temp2_1 > 60 || p.temp2_2 > 60 || p.temp2_3 > 60 || p.temp2_4 > 60
        attrs.heat = true

    summary = Miner.summary()
    if ifError(summary)
        Actor.setState
            'new_state': 'ERROR'
        return

    p = JSON.parse(summary.result)
    attrs["ghs_av"] = p["GHS av"] 
    attrs["hardware_errors"] = p["Hardware Errors"] 
   
    Actor.setState
        new_state: 'ENABLED'
        attribute_values: attrs

entityAction = (entityId, actionName)->
    switch actionName
        when 'CHECK' then checkStatus()
`

	const L3PlusStatsJson = `{"STATUS":[{"STATUS":"S","When":1633268527,"Code":70,"Msg":"CGMiner stats","Description":"cgminer 4.9.0 rwglr"}],"STATS":[{"CGMiner":"4.9.0 rwglr","Miner":"1.0.1.3","CompileTime":"Fri 11 Jun 15:32:42 MSK 2021","Type":"Antminer L3+ (024)"}{"STATS":0,"ID":"L30","Elapsed":373469,"Calls":0,"Wait":0.000000,"Max":0.000000,"Min":99999999.000000,"GHS 5s":"415.065","GHS av":411.84,"miner_count":3,"frequency":"425","fan_num":2,"frequency1":"425","frequency2":"422","frequency3":"423","frequency4":"419","fan1":3180,"fan2":3090,"temp_num":3,"temp1":0,"temp2":37,"temp3":37,"temp4":36,"temp2_1":0,"temp2_2":46,"temp2_3":44,"temp2_4":42,"temp31":0,"temp32":0,"temp33":0,"temp34":0,"temp4_1":0,"temp4_2":0,"temp4_3":0,"temp4_4":0,"temp_max":38,"Device Hardware%":0.0000,"no_matching_work":327,"chain_acn1":0,"chain_acn2":72,"chain_acn3":72,"chain_acn4":72,"chain_acs1":"","chain_acs2":" oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo","chain_acs3":" oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo","chain_acs4":" oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo","chain_hw1":0,"chain_hw2":26,"chain_hw3":43,"chain_hw4":258,"chain_rate1":"","chain_rate2":"138.50","chain_rate3":"138.85","chain_rate4":"137.72","chain_power1":0.00,"chain_power2":187.79,"chain_power3":188.24,"chain_power4":186.46,"chain_power":562.48,"voltage1":9.44,"voltage2":9.44,"voltage3":9.44,"voltage4":9.44}],"id":1}`
	const L3PlusSummaryJson = `{"STATUS":[{"STATUS":"S","When":1633279121,"Code":11,"Msg":"Summary","Description":"cgminer 4.9.0 rwglr"}],"SUMMARY":[{"Elapsed":384063,"GHS 5s":"412.925","GHS av":411.85,"Found Blocks":0,"Getworks":37619,"Accepted":46423,"Rejected":530,"Hardware Errors":338,"Utility":7.25,"Discarded":226550,"Stale":21,"Get Failures":169,"Local Work":338070,"Remote Failures":4,"Network Blocks":2603,"Total MH":158174645.0000,"Work Utility":377447.21,"Difficulty Accepted":2397855744.00000000,"Difficulty Rejected":18202624.00000000,"Difficulty Stale":0.00000000,"Best Share":7510923133,"Device Hardware%":0.0000,"Device Rejected%":0.7534,"Pool Rejected%":0.7534,"Pool Stale%":0.0000,"Last getwork":1633279119}],"id":1}`

	var getFixture = func(str string) []byte {
		return append([]byte(str), byte(0x00))
	}

	Convey("cgminer", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService scripts.ScriptService,
			entityManager entity_manager.EntityManager,
			mqttServer mqtt.MqttServ,
			eventBus event_bus.EventBus,
			pluginManager common.PluginManager) {

			eventBus.Purge()
			scriptService.Purge()

			err := migrations.Purge()
			ctx.So(err, ShouldBeNil)

			// bind convey
			RegisterConvey(scriptService, ctx)

			// register plugins
			err = AddPlugin(adaptors, "cgminer")
			ctx.So(err, ShouldBeNil)

			go mqttServer.Start()

			// add scripts
			// ------------------------------------------------
			plugScript := &m.Script{
				Lang:        common.ScriptLangCoffee,
				Name:        "plug script",
				Source:      cgminerSourceScript,
				Description: "cgminer script",
			}

			engineScript, err := scriptService.NewEngine(plugScript)
			So(err, ShouldBeNil)
			err = engineScript.Compile()
			So(err, ShouldBeNil)

			plugScript.Id, err = adaptors.Script.Add(plugScript)
			So(err, ShouldBeNil)

			// add entity
			// ------------------------------------------------
			var port = GetPort()
			var host = "127.0.0.1"

			l3Ent := GetNewBitmineL3("device1")
			l3Ent.Settings[cgminer.SettingHost].Value = host
			l3Ent.Settings[cgminer.SettingPort].Value = port
			l3Ent.Actions = []*m.EntityAction{
				{
					Name:        "ENABLE",
					Description: "включить",
					Script:      plugScript,
				},
				{
					Name:        "DISABLE",
					Description: "выключить",
					Script:      plugScript,
				},
				{
					Name:        "CHECK",
					Description: "condition check",
					Script:      plugScript,
				},
			}
			l3Ent.States = []*m.EntityState{
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
			}
			l3Ent.Attributes = m.Attributes{
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
			err = adaptors.Entity.Add(l3Ent)
			ctx.So(err, ShouldBeNil)
			_, err = adaptors.EntityStorage.Add(m.EntityStorage{
				EntityId:   l3Ent.Id,
				Attributes: l3Ent.Attributes.Serialize(),
			})
			So(err, ShouldBeNil)

			// ------------------------------------------------
			pluginManager.Start()
			entityManager.LoadEntities(pluginManager)

			defer func() {
				mqttServer.Shutdown()
				entityManager.Shutdown()
				pluginManager.Shutdown()
			}()

			time.Sleep(time.Millisecond * 500)

			t.Run("antminer L3+", func(t *testing.T) {
				Convey("check status", t, func(ctx C) {

					ctx2, cancel := context2.WithCancel(context2.Background())
					go MockTCPServer(ctx2, host, port, getFixture(L3PlusStatsJson), getFixture(L3PlusSummaryJson))
					time.Sleep(time.Millisecond * 500)
					entityManager.CallAction(l3Ent.Id, "CHECK", nil)
					time.Sleep(time.Second)
					cancel()

					lastState, err := adaptors.EntityStorage.GetLastByEntityId(l3Ent.Id)
					ctx.So(err, ShouldBeNil)

					ctx.So(lastState.Attributes["chain1_temp_chip"], ShouldEqual, 0)
					ctx.So(lastState.Attributes["chain1_temp_pcb"], ShouldEqual, 0)
					ctx.So(lastState.Attributes["chain2_temp_chip"], ShouldEqual, 46)
					ctx.So(lastState.Attributes["chain2_temp_pcb"], ShouldEqual, 37)
					ctx.So(lastState.Attributes["chain3_temp_chip"], ShouldEqual, 44)
					ctx.So(lastState.Attributes["chain3_temp_pcb"], ShouldEqual, 37)
					ctx.So(lastState.Attributes["chain4_temp_chip"], ShouldEqual, 42)
					ctx.So(lastState.Attributes["chain4_temp_pcb"], ShouldEqual, 36)
					ctx.So(lastState.Attributes["fan1"], ShouldEqual, 3180)
					ctx.So(lastState.Attributes["fan2"], ShouldEqual, 3090)
					ctx.So(lastState.Attributes["ghs_av"], ShouldEqual, 411.85)
					ctx.So(lastState.Attributes["hardware_errors"], ShouldEqual, 338)
					ctx.So(lastState.Attributes["heat"], ShouldEqual, false)
				})
			})

		})
	})
}
