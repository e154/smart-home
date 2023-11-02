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

package cgminer

import (
	"context"
	"github.com/e154/smart-home/common/debug"
	"testing"
	"time"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/cgminer"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
	. "github.com/e154/smart-home/tests/plugins"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCgMiner(t *testing.T) {

	const cgminerSourceScript = `

ifError =(res)->
    return !res || res.error || res.Error

checkStatus =->
    stats = Miner.stats()
    if ifError(stats)
        EntitySetStateName ENTITY_ID, 'ERROR'
        return
    p = unmarshal stats.result
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

    EntitySetState ENTITY_ID,
        new_state: 'ENABLED'
        attribute_values: attrs
        storage_save: true

checkSum =->
    summary = Miner.summary()
    if ifError(summary)
        EntitySetStateName ENTITY_ID, 'ERROR'
        return
    p = unmarshal summary.result
    attrs = {}
    attrs["ghs_av"] = p["GHS av"] 
    attrs["hardware_errors"] = p["Hardware Errors"] 
    EntitySetState ENTITY_ID,
        new_state: 'ENABLED'
        attribute_values: attrs
        storage_save: true

checkDevs =->
    devs = Miner.devs()
    if ifError(devs)
        EntitySetStateName ENTITY_ID, 'ERROR'
        return
    p = unmarshal devs.result
    attrs = {}
    attrs["ghs_av"] = p["GHS av"] 
    attrs["hardware_errors"] = p["Hardware Errors"] 
    EntitySetState ENTITY_ID,
        new_state: 'ENABLED'
        attribute_values: attrs
        storage_save: true

checkPools =->
    pools = Miner.pools()
    if ifError(pools)
        EntitySetStateName ENTITY_ID, 'ERROR'
        return
    p = unmarshal pools.result
    So(p.length, 'ShouldEqual', '4')

    So(p[0].POOL, 'ShouldEqual', '0')
    So(p[0].Accepted, 'ShouldEqual', '47169')
    So(p[0]["Stratum URL"], 'ShouldEqual', 'litecoinpool.org')
    So(p[0]["User"], 'ShouldEqual', 'user.4')
    So(p[0]["URL"], 'ShouldEqual', 'stratum+tcp://litecoinpool.org:3333')

    So(p[1].POOL, 'ShouldEqual', '1')
    So(p[1].Accepted, 'ShouldEqual', '0')
    So(p[1]["Stratum URL"], 'ShouldEqual', '')
    So(p[1]["User"], 'ShouldEqual', 'user.4')
    So(p[1]["URL"], 'ShouldEqual', 'stratum+tcp://us.litecoinpool.org:3333')

    So(p[2].POOL, 'ShouldEqual', '2')
    So(p[2].Accepted, 'ShouldEqual', '0')
    So(p[2]["Stratum URL"], 'ShouldEqual', '')
    So(p[2]["User"], 'ShouldEqual', 'user.4')
    So(p[2]["URL"], 'ShouldEqual', 'stratum+tcp://us2.litecoinpool.org:3333')

    So(p[3].POOL, 'ShouldEqual', '3')
    So(p[3].Accepted, 'ShouldEqual', '5502')
    So(p[3]["Stratum URL"], 'ShouldEqual', 'proto3.bytewall.net')
    So(p[3]["User"], 'ShouldEqual', '*')
    So(p[3]["URL"], 'ShouldEqual', '*')

checkVer =(entityId)->
    ver = Miner.version()
    if ifError(ver)
        EntitySetStateName ENTITY_ID, 'ERROR'
        return
    p = unmarshal ver.result
    So(p.API, 'ShouldEqual', '3.1')
    So(p.Miner, 'ShouldEqual', '1.0.1.3')
    So(p.CompileTime, 'ShouldEqual', 'Fri 11 Jun 15:32:42 MSK 2021')
    So(p.Type, 'ShouldEqual', 'Antminer L3+ (024)')

entityAction = (entityId, actionName)->
    switch actionName
        when 'CHECK' then checkStatus()
        when 'SUM' then checkSum()
        when 'DEVS' then checkDevs()
        when 'POOLS' then checkPools()
        when 'VER' then checkVer(entityId)
`

	const L3PlusStatsJson = `{"STATUS":[{"STATUS":"S","When":1633268527,"Code":70,"Msg":"CGMiner stats","Description":"cgminer 4.9.0 rwglr"}],"STATS":[{"CGMiner":"4.9.0 rwglr","Miner":"1.0.1.3","CompileTime":"Fri 11 Jun 15:32:42 MSK 2021","Type":"Antminer L3+ (024)"}{"STATS":0,"ID":"L30","Elapsed":373469,"Calls":0,"Wait":0.000000,"Max":0.000000,"Min":99999999.000000,"GHS 5s":"415.065","GHS av":411.84,"miner_count":3,"frequency":"425","fan_num":2,"frequency1":"425","frequency2":"422","frequency3":"423","frequency4":"419","fan1":3180,"fan2":3090,"temp_num":3,"temp1":0,"temp2":37,"temp3":37,"temp4":36,"temp2_1":0,"temp2_2":46,"temp2_3":44,"temp2_4":42,"temp31":0,"temp32":0,"temp33":0,"temp34":0,"temp4_1":0,"temp4_2":0,"temp4_3":0,"temp4_4":0,"temp_max":38,"Device Hardware%":0.0000,"no_matching_work":327,"chain_acn1":0,"chain_acn2":72,"chain_acn3":72,"chain_acn4":72,"chain_acs1":"","chain_acs2":" oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo","chain_acs3":" oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo","chain_acs4":" oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo","chain_hw1":0,"chain_hw2":26,"chain_hw3":43,"chain_hw4":258,"chain_rate1":"","chain_rate2":"138.50","chain_rate3":"138.85","chain_rate4":"137.72","chain_power1":0.00,"chain_power2":187.79,"chain_power3":188.24,"chain_power4":186.46,"chain_power":562.48,"voltage1":9.44,"voltage2":9.44,"voltage3":9.44,"voltage4":9.44}],"id":1}`
	const L3PlusSummaryJson = `{"STATUS":[{"STATUS":"S","When":1633279121,"Code":11,"Msg":"Summary","Description":"cgminer 4.9.0 rwglr"}],"SUMMARY":[{"Elapsed":384063,"GHS 5s":"412.925","GHS av":411.85,"Found Blocks":0,"Getworks":37619,"Accepted":46423,"Rejected":530,"Hardware Errors":338,"Utility":7.25,"Discarded":226550,"Stale":21,"Get Failures":169,"Local Work":338070,"Remote Failures":4,"Network Blocks":2603,"Total MH":158174645.0000,"Work Utility":377447.21,"Difficulty Accepted":2397855744.00000000,"Difficulty Rejected":18202624.00000000,"Difficulty Stale":0.00000000,"Best Share":7510923133,"Device Hardware%":0.0000,"Device Rejected%":0.7534,"Pool Rejected%":0.7534,"Pool Stale%":0.0000,"Last getwork":1633279119}],"id":1}`
	const L3PlusPoolsJson = `{"STATUS":[{"STATUS":"S","When":1633329118,"Code":7,"Msg":"4 Pool(s)","Description":"cgminer 4.9.0 rwglr"}],"POOLS":[{"POOL":0,"URL":"stratum+tcp://litecoinpool.org:3333","Status":"Alive","Priority":0,"Quota":1,"Long Poll":"N","Getworks":19515,"Accepted":47169,"Rejected":331,"Discarded":250226,"Stale":16,"Get Failures":2,"Remote Failures":0,"User":"user.4","Last Share Time":"0:00:06","Diff":"65.5K","Diff1 Shares":10494726,"Proxy Type":"","Proxy":"","Difficulty Accepted":2659909632.00000000,"Difficulty Rejected":18219008.00000000,"Difficulty Stale":0.00000000,"Last Share Difficulty":65536.00000000,"Has Stratum":true,"Stratum Active":true,"Stratum URL":"litecoinpool.org","Has GBT":false,"Best Share":7510923133,"Pool Rejected%":0.6803,"Pool Stale%":0.0000},{"POOL":1,"URL":"stratum+tcp://us.litecoinpool.org:3333","Status":"Alive","Priority":1,"Quota":1,"Long Poll":"N","Getworks":1,"Accepted":0,"Rejected":0,"Discarded":0,"Stale":0,"Get Failures":0,"Remote Failures":0,"User":"user.4","Last Share Time":"0","Diff":"","Diff1 Shares":0,"Proxy Type":"","Proxy":"","Difficulty Accepted":0.00000000,"Difficulty Rejected":0.00000000,"Difficulty Stale":0.00000000,"Last Share Difficulty":0.00000000,"Has Stratum":true,"Stratum Active":false,"Stratum URL":"","Has GBT":false,"Best Share":0,"Pool Rejected%":0.0000,"Pool Stale%":0.0000},{"POOL":2,"URL":"stratum+tcp://us2.litecoinpool.org:3333","Status":"Alive","Priority":2,"Quota":1,"Long Poll":"N","Getworks":1,"Accepted":0,"Rejected":0,"Discarded":0,"Stale":0,"Get Failures":0,"Remote Failures":0,"User":"user.4","Last Share Time":"0","Diff":"","Diff1 Shares":0,"Proxy Type":"","Proxy":"","Difficulty Accepted":0.00000000,"Difficulty Rejected":0.00000000,"Difficulty Stale":0.00000000,"Last Share Difficulty":0.00000000,"Has Stratum":true,"Stratum Active":false,"Stratum URL":"","Has GBT":false,"Best Share":0,"Pool Rejected%":0.0000,"Pool Stale%":0.0000},{"POOL":3,"URL":"*","Status":"Alive","Priority":999,"Quota":1,"Long Poll":"N","Getworks":22556,"Accepted":5502,"Rejected":288,"Discarded":5377,"Stale":5,"Get Failures":200,"Remote Failures":4,"User":"*","Last Share Time":"0:14:56","Diff":"131K","Diff1 Shares":225056,"Proxy Type":"","Proxy":"","Difficulty Accepted":54272000.00000000,"Difficulty Rejected":2957312.00000000,"Difficulty Stale":0.00000000,"Last Share Difficulty":131072.00000000,"Has Stratum":true,"Stratum Active":true,"Stratum URL":"proto3.bytewall.net","Has GBT":false,"Best Share":66992548,"Pool Rejected%":5.1675,"Pool Stale%":0.0000}],"id":1}`
	const L3PlusVerJson = `{"STATUS":[{"STATUS":"S","When":1633332261,"Code":22,"Msg":"CGMiner versions","Description":"cgminer 4.9.0 rwglr"}],"VERSION":[{"CGMiner":"4.9.0 rwglr","API":"3.1","Miner":"1.0.1.3","CompileTime":"Fri 11 Jun 15:32:42 MSK 2021","Type":"Antminer L3+ (024)"}],"id":1}`

	var getFixture = func(str string) []byte {
		return append([]byte(str), byte(0x00))
	}

	Convey("cgminer", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			scriptService scripts.ScriptService,
			supervisor supervisor.Supervisor,
			eventBus bus.Bus) {

			// register plugins
			AddPlugin(adaptors, "cgminer")

			supervisor.Start(context.Background())
			WaitSupervisor(eventBus)

			// bind convey
			RegisterConvey(scriptService, ctx)

			// add scripts
			// ------------------------------------------------

			plugScript, err := AddScript("cgminer script", cgminerSourceScript, adaptors, scriptService)
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
				{
					Name:        "SUM",
					Description: "summary",
					Script:      plugScript,
				},
				{
					Name:        "DEVS",
					Description: "devs",
					Script:      plugScript,
				},
				{
					Name:        "POOLS",
					Description: "pools",
					Script:      plugScript,
				},
				{
					Name:        "VER",
					Description: "ver",
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
			err = adaptors.Entity.Add(context.Background(), l3Ent)
			ctx.So(err, ShouldBeNil)
			_, err = adaptors.EntityStorage.Add(context.Background(), &m.EntityStorage{
				EntityId:   l3Ent.Id,
				Attributes: l3Ent.Attributes.Serialize(),
			})
			So(err, ShouldBeNil)

			eventBus.Publish("system/models/entities/"+l3Ent.Id.String(), events.EventCreatedEntityModel{
				EntityId: l3Ent.Id,
			})

			time.Sleep(time.Second)

			// ------------------------------------------------

			t.Run("antminer L3+ stats", func(t *testing.T) {
				Convey("stats", t, func(ctx C) {

					ctx2, cancel := context.WithCancel(context.Background())
					go func() { _ = MockTCPServer(ctx2, host, port, getFixture(L3PlusStatsJson)) }()
					time.Sleep(time.Millisecond * 500)
					supervisor.CallAction(l3Ent.Id, "CHECK", nil)
					time.Sleep(time.Second)
					cancel()

					lastState, err := adaptors.EntityStorage.GetLastByEntityId(context.Background(), l3Ent.Id)
					ctx.So(err, ShouldBeNil)

					debug.Println(lastState)

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
				})
			})

			t.Run("antminer L3+ summary", func(t *testing.T) {
				Convey("summary", t, func(ctx C) {

					ctx2, cancel := context.WithCancel(context.Background())
					go func() { _ = MockTCPServer(ctx2, host, port, getFixture(L3PlusSummaryJson)) }()
					time.Sleep(time.Millisecond * 500)
					supervisor.CallAction(l3Ent.Id, "SUM", nil)
					time.Sleep(time.Second)
					cancel()

					lastState, err := adaptors.EntityStorage.GetLastByEntityId(context.Background(), l3Ent.Id)
					ctx.So(err, ShouldBeNil)

					ctx.So(lastState.Attributes["ghs_av"], ShouldEqual, 411.85)
					ctx.So(lastState.Attributes["hardware_errors"], ShouldEqual, 338)
				})
			})

			t.Run("antminer L3+ pools", func(t *testing.T) {
				Convey("pools", t, func(ctx C) {

					ctx2, cancel := context.WithCancel(context.Background())
					go func() { _ = MockTCPServer(ctx2, host, port, getFixture(L3PlusPoolsJson)) }()
					time.Sleep(time.Millisecond * 500)
					supervisor.CallAction(l3Ent.Id, "POOLS", nil)
					time.Sleep(time.Second)
					cancel()
				})
			})

			t.Run("antminer L3+ version", func(t *testing.T) {
				Convey("version", t, func(ctx C) {

					ctx2, cancel := context.WithCancel(context.Background())
					go func() { _ = MockTCPServer(ctx2, host, port, getFixture(L3PlusVerJson)) }()
					time.Sleep(time.Millisecond * 500)
					supervisor.CallAction(l3Ent.Id, "VER", nil)
					time.Sleep(time.Second)
					cancel()
				})
			})

		})
	})
}
