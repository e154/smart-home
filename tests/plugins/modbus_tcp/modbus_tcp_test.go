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

package modbus_tcp

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/e154/smart-home/common/events"

	"github.com/e154/bus"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/modbus_tcp"
	"github.com/e154/smart-home/plugins/node"
	"github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
	. "github.com/e154/smart-home/tests/plugins"
	. "github.com/smartystreets/goconvey/convey"
)

func TestModbusTcp(t *testing.T) {

	const plugActionOnOffSourceScript = `

writeRegisters =(d, c, r)->
    return ModbusTcp 'WriteMultipleRegisters', d, c, r

checkStatus =->
    COMMAND = []
    FUNC = 'ReadHoldingRegisters'
    ADDRESS = 0
    COUNT = 16

    res = ModbusTcp FUNC, ADDRESS, COUNT, COMMAND
    So(res.error, 'ShouldEqual', '')
    So(res.result, 'ShouldEqual', '[1 0 1]')
    So(res.time, 'ShouldNotBeBlank', '')

doOnAction = ->
    res = writeRegisters(0, 1, [1])
    So(res.error, 'ShouldEqual', '')
    So(res.result, 'ShouldEqual', '[]')
    So(res.time, 'ShouldNotBeBlank', '')

doOnErrAction = ->
    res = writeRegisters(0, 1, [1])
    So(res.error, 'ShouldEqual', 'some error')
    So(res.result, 'ShouldEqual', '[]')
    So(res.time, 'ShouldNotBeBlank', '')

doOffAction = ->
    res = writeRegisters(0, 1, [0])
    So(res.error, 'ShouldEqual', '')
    So(res.result, 'ShouldEqual', '[]')
    So(res.time, 'ShouldNotBeBlank', '')

entityAction = (entityId, actionName)->
    #print '---action on/off--'
    switch actionName
        when 'ON' then doOnAction()
        when 'OFF' then doOffAction()
        when 'CHECK' then checkStatus()
        when 'ON_WITH_ERR' then doOnErrAction()
`

	Convey("modbus_tcp", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService scripts.ScriptService,
			supervisor supervisor.Supervisor,
			mqttServer mqtt.MqttServ,
			automation automation.Automation,
			eventBus bus.Bus) {

			// register plugins
			AddPlugin(adaptors, "node")
			AddPlugin(adaptors, "triggers")
			AddPlugin(adaptors, "modbus_tcp")

			serviceCh := WaitService(eventBus, time.Second*5, "Supervisor", "Automation", "Mqtt")
			pluginsCh := WaitPlugins(eventBus, time.Second*5, "triggers", "modbus_tcp", "node")
			automation.Start()
			go mqttServer.Start()
			supervisor.Start(context.Background())
			defer automation.Shutdown()
			defer mqttServer.Shutdown()
			defer supervisor.Shutdown(context.Background())
			So(<-serviceCh, ShouldBeTrue)
			So(<-pluginsCh, ShouldBeTrue)

			// bind convey
			RegisterConvey(scriptService, ctx)

			// add scripts
			// ------------------------------------------------

			plugActionOnOffScript, err := AddScript("plug script", plugActionOnOffSourceScript, adaptors, scriptService)
			So(err, ShouldBeNil)

			// add entity
			// ------------------------------------------------
			nodeEnt := GetNewNode("second")
			err = adaptors.Entity.Add(context.Background(), nodeEnt)
			So(err, ShouldBeNil)

			plugEnt := GetNewModbusTcp("plug")
			plugEnt.ParentId = &nodeEnt.Id
			plugEnt.Actions = []*m.EntityAction{
				{
					Name:        "ON",
					Description: "включить",
					Script:      plugActionOnOffScript,
				},
				{
					Name:        "OFF",
					Description: "выключить",
					Script:      plugActionOnOffScript,
				},
				{
					Name:        "CHECK",
					Description: "condition check",
					Script:      plugActionOnOffScript,
				},
				{
					Name:        "ON_WITH_ERR",
					Description: "error case",
					Script:      plugActionOnOffScript,
				},
			}
			plugEnt.States = []*m.EntityState{
				{
					Name:        "ON",
					Description: "on state",
				},
				{
					Name:        "OFF",
					Description: "off state",
				},
				{
					Name:        "ERROR",
					Description: "error state",
				},
			}
			plugEnt.Settings[modbus_tcp.AttrSlaveId].Value = 1
			plugEnt.Settings[modbus_tcp.AttrAddressPort].Value = "office:502"
			err = adaptors.Entity.Add(context.Background(), plugEnt)
			So(err, ShouldBeNil)
			_, err = adaptors.EntityStorage.Add(context.Background(), &m.EntityStorage{
				EntityId:   plugEnt.Id,
				Attributes: plugEnt.Attributes.Serialize(),
			})
			So(err, ShouldBeNil)

			eventBus.Publish("system/models/entities/"+nodeEnt.Id.String(), events.EventCreatedEntityModel{
				EntityId: nodeEnt.Id,
			})
			eventBus.Publish("system/models/entities/"+nodeEnt.Id.String(), events.EventCreatedEntityModel{
				EntityId: plugEnt.Id,
			})

			time.Sleep(time.Second)

			// ------------------------------------------------

			ch := make(chan []byte)
			mqttCli := mqttServer.NewClient("cli")
			_ = mqttCli.Subscribe("system/plugins/node/second/req/#", func(cli mqtt.MqttCli, message mqtt.Message) {
				ch <- message.Payload
			})
			defer mqttCli.UnsubscribeAll()

			// commands
			t.Run("on command", func(t *testing.T) {
				Convey("stats", t, func(ctx C) {
					supervisor.CallAction(plugEnt.Id, "ON", nil)

					// wait message
					req, ok := WaitT[[]byte](time.Second*2, ch)
					ctx.So(ok, ShouldBeTrue)

					// what see node
					request := node.MessageRequest{}
					err = json.Unmarshal(req, &request)
					ctx.So(err, ShouldBeNil)

					cmd := modbus_tcp.ModBusCommand{}
					err = json.Unmarshal(request.Command, &cmd)
					ctx.So(err, ShouldBeNil)

					prop := map[string]interface{}{}
					err = json.Unmarshal(request.Properties, &prop)
					ctx.So(err, ShouldBeNil)

					ctx.So(request.EntityId, ShouldEqual, plugEnt.Id)
					ctx.So(request.DeviceType, ShouldEqual, modbus_tcp.DeviceTypeModbusTcp)
					ctx.So(prop["slave_id"], ShouldEqual, 1)
					ctx.So(prop["address_port"], ShouldEqual, "office:502")

					ctx.So(cmd.Function, ShouldEqual, "WriteMultipleRegisters")
					ctx.So(cmd.Address, ShouldEqual, 0)
					ctx.So(cmd.Count, ShouldEqual, 1)
					ctx.So(cmd.Command, ShouldResemble, []uint16{1})

					// response from node
					r := modbus_tcp.ModBusResponse{
						Error:  "",
						Result: []uint16{},
					}
					b, _ := json.Marshal(r)
					resp := node.MessageResponse{
						EntityId:   plugEnt.Id,
						DeviceType: modbus_tcp.DeviceTypeModbusTcp,
						Properties: nil,
						Response:   b,
						Status:     "",
					}
					b, _ = json.Marshal(resp)
					_ = mqttCli.Publish("system/plugins/node/second/resp/plugin.test", b)
					_ = mqttCli.Publish(fmt.Sprintf("system/plugins/node/second/resp/%s", plugEnt.Id), b)

					time.Sleep(time.Millisecond * 500)
				})
			})

			t.Run("off command", func(t *testing.T) {
				Convey("stats", t, func(ctx C) {
					supervisor.CallAction(plugEnt.Id, "OFF", nil)

					// wait message
					req, ok := WaitT[[]byte](time.Second*2, ch)
					ctx.So(ok, ShouldBeTrue)

					// what see node
					request := node.MessageRequest{}
					err = json.Unmarshal(req, &request)
					ctx.So(err, ShouldBeNil)

					cmd := modbus_tcp.ModBusCommand{}
					err = json.Unmarshal(request.Command, &cmd)
					ctx.So(err, ShouldBeNil)

					prop := map[string]interface{}{}
					err = json.Unmarshal(request.Properties, &prop)
					ctx.So(err, ShouldBeNil)

					ctx.So(request.EntityId, ShouldEqual, plugEnt.Id)
					ctx.So(request.DeviceType, ShouldEqual, modbus_tcp.DeviceTypeModbusTcp)
					ctx.So(prop["slave_id"], ShouldEqual, 1)
					ctx.So(prop["address_port"], ShouldEqual, "office:502")

					ctx.So(cmd.Function, ShouldEqual, "WriteMultipleRegisters")
					ctx.So(cmd.Address, ShouldEqual, 0)
					ctx.So(cmd.Count, ShouldEqual, 1)
					ctx.So(cmd.Command, ShouldResemble, []uint16{0})

					// response from node
					r := modbus_tcp.ModBusResponse{
						Error:  "",
						Result: []uint16{},
					}
					b, _ := json.Marshal(r)
					resp := node.MessageResponse{
						EntityId:   plugEnt.Id,
						DeviceType: modbus_tcp.DeviceTypeModbusTcp,
						Properties: nil,
						Response:   b,
						Status:     "",
					}
					b, _ = json.Marshal(resp)
					_ = mqttCli.Publish("system/plugins/node/second/resp/plugin.test", b)
					_ = mqttCli.Publish(fmt.Sprintf("system/plugins/node/second/resp/%s", plugEnt.Id), b)

					time.Sleep(time.Millisecond * 500)
				})
			})

			t.Run("check command", func(t *testing.T) {
				Convey("stats", t, func(ctx C) {
					supervisor.CallAction(plugEnt.Id, "CHECK", nil)

					// wait message
					req, ok := WaitT[[]byte](time.Second*2, ch)
					ctx.So(ok, ShouldBeTrue)

					// what see node
					request := node.MessageRequest{}
					err = json.Unmarshal(req, &request)
					ctx.So(err, ShouldBeNil)

					cmd := modbus_tcp.ModBusCommand{}
					err = json.Unmarshal(request.Command, &cmd)
					ctx.So(err, ShouldBeNil)

					prop := map[string]interface{}{}
					err = json.Unmarshal(request.Properties, &prop)
					ctx.So(err, ShouldBeNil)

					ctx.So(request.EntityId, ShouldEqual, plugEnt.Id)
					ctx.So(request.DeviceType, ShouldEqual, modbus_tcp.DeviceTypeModbusTcp)
					ctx.So(prop["slave_id"], ShouldEqual, 1)
					ctx.So(prop["address_port"], ShouldEqual, "office:502")

					ctx.So(cmd.Function, ShouldEqual, "ReadHoldingRegisters")
					ctx.So(cmd.Address, ShouldEqual, 0)
					ctx.So(cmd.Count, ShouldEqual, 16)
					ctx.So(cmd.Command, ShouldResemble, []uint16{})

					// response from node
					r := modbus_tcp.ModBusResponse{
						Error:  "",
						Result: []uint16{1, 0, 1},
					}
					b, _ := json.Marshal(r)
					resp := node.MessageResponse{
						EntityId:   plugEnt.Id,
						DeviceType: modbus_tcp.DeviceTypeModbusTcp,
						Properties: nil,
						Response:   b,
						Status:     "",
					}
					b, _ = json.Marshal(resp)
					_ = mqttCli.Publish("system/plugins/node/second/resp/plugin.test", b)
					_ = mqttCli.Publish(fmt.Sprintf("system/plugins/node/second/resp/%s", plugEnt.Id), b)

					time.Sleep(time.Millisecond * 500)
				})
			})

			t.Run("bad command", func(t *testing.T) {
				Convey("stats", t, func(ctx C) {
					supervisor.CallAction(plugEnt.Id, "NULL", nil)

					// wait message
					_, ok := WaitT[[]byte](time.Second*2, ch)
					ctx.So(ok, ShouldBeFalse)
				})
			})

			t.Run("response with error", func(t *testing.T) {
				Convey("stats", t, func(ctx C) {
					supervisor.CallAction(plugEnt.Id, "ON_WITH_ERR", nil)

					// wait message
					_, ok := WaitT[[]byte](time.Second*2, ch)
					ctx.So(ok, ShouldBeTrue)

					r := modbus_tcp.ModBusResponse{
						Error: "some error",
					}
					b, _ := json.Marshal(r)
					resp := node.MessageResponse{
						EntityId:   plugEnt.Id,
						DeviceType: modbus_tcp.DeviceTypeModbusTcp,
						Properties: nil,
						Response:   b,
						Status:     "",
					}
					b, _ = json.Marshal(resp)
					_ = mqttCli.Publish("system/plugins/node/second/resp/plugin.test", b)
					_ = mqttCli.Publish(fmt.Sprintf("system/plugins/node/second/resp/%s", plugEnt.Id), b)

				})
			})
		})
	})
}
