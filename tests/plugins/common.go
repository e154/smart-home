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

package plugins

import (
	"context"
	"fmt"
	"go.uber.org/atomic"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/phayes/freeport"
	"github.com/smartystreets/goconvey/convey"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/alexa"
	"github.com/e154/smart-home/plugins/cgminer"
	"github.com/e154/smart-home/plugins/cgminer/bitmine"
	"github.com/e154/smart-home/plugins/modbus_rtu"
	"github.com/e154/smart-home/plugins/modbus_tcp"
	"github.com/e154/smart-home/plugins/moon"
	"github.com/e154/smart-home/plugins/node"
	"github.com/e154/smart-home/plugins/scene"
	"github.com/e154/smart-home/plugins/sun"
	"github.com/e154/smart-home/plugins/telegram"
	"github.com/e154/smart-home/plugins/weather"
	"github.com/e154/smart-home/plugins/zigbee2mqtt"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/scripts"
)

// GetNewButton ...
func GetNewButton(id string, scripts []*m.Script) *m.Entity {
	return &m.Entity{
		Id:           common.EntityId(id),
		Description:  "MiJia wireless switch",
		PluginName:   zigbee2mqtt.EntityZigbee2mqtt,
		Scripts:      scripts,
		AutoLoad:     true,
		RestoreState: true,
		Attributes: m.Attributes{
			"click": &m.Attribute{
				Name: "click",
				Type: common.AttributeString,
			},
			"action": &m.Attribute{
				Name: "action",
				Type: common.AttributeString,
			},
			"battery": &m.Attribute{
				Name: "battery",
				Type: common.AttributeInt,
			},
			"voltage": &m.Attribute{
				Name: "voltage",
				Type: common.AttributeInt,
			},
			"linkquality": &m.Attribute{
				Name: "linkquality",
				Type: common.AttributeInt,
			},
		},
		States: []*m.EntityState{
			{
				Name:        "LONG_CLICK",
				Description: "long click",
			},
			{
				Name:        "LONG_ACTION",
				Description: "long action",
			},
			{
				Name:        "SINGLE_CLICK",
				Description: "single click",
			},
			{
				Name:        "SINGLE_ACTION",
				Description: "single action",
			},
			{
				Name:        "DOUBLE_CLICK",
				Description: "double click",
			},
			{
				Name:        "DOUBLE_ACTION",
				Description: "double action",
			},
			{
				Name:        "TRIPLE_CLICK",
				Description: "triple click",
			},
			{
				Name:        "TRIPLE_ACTION",
				Description: "triple action",
			},
			{
				Name:        "QUADRUPLE_CLICK",
				Description: "quadruple click",
			},
			{
				Name:        "QUADRUPLE_ACTION",
				Description: "quadruple action",
			},
			{
				Name:        "MANY_CLICK",
				Description: "many click",
			},
			{
				Name:        "MANY_ACTION",
				Description: "many action",
			},
			{
				Name:        "LONG_RELEASE_CLICK",
				Description: "long_release click",
			},
			{
				Name:        "HOLD_ACTION",
				Description: "hold action",
			},
			{
				Name:        "RELEASE_ACTION",
				Description: "release action",
			},
		},
	}
}

// GetNewPlug ...
func GetNewPlug(id string, scrits []*m.Script) *m.Entity {
	return &m.Entity{
		Id:           common.EntityId(id),
		Description:  "MiJia power plug ZigBee",
		PluginName:   zigbee2mqtt.EntityZigbee2mqtt,
		Scripts:      scrits,
		AutoLoad:     true,
		RestoreState: true,
		Attributes: m.Attributes{
			"power": &m.Attribute{
				Name: "power",
				Type: common.AttributeInt,
			},
			"state": &m.Attribute{
				Name: "state",
				Type: common.AttributeString,
			},
			"voltage": &m.Attribute{
				Name: "voltage",
				Type: common.AttributeInt,
			},
			"consumption": &m.Attribute{
				Name: "consumption",
				Type: common.AttributeString,
			},
			"linkquality": &m.Attribute{
				Name: "linkquality",
				Type: common.AttributeInt,
			},
			"temperature": &m.Attribute{
				Name: "temperature",
				Type: common.AttributeInt,
			},
		},
		States: []*m.EntityState{
			{
				Name:        "ON",
				Description: "on state",
			},
			{
				Name:        "OFF",
				Description: "off state",
			},
		},
	}
}

// GetNewScene ...
func GetNewScene(id string, scripts []*m.Script) *m.Entity {
	return &m.Entity{
		Id:           common.EntityId(id),
		Description:  "scene",
		PluginName:   scene.EntityScene,
		Scripts:      scripts,
		AutoLoad:     true,
		RestoreState: true,
	}
}

// GetNewNode ...
func GetNewNode(name string) *m.Entity {
	settings := node.NewSettings()
	settings[node.AttrNodeLogin].Value = "node1"
	settings[node.AttrNodePass].Value = "node1"
	return &m.Entity{
		Id:           common.EntityId(fmt.Sprintf("node.%s", name)),
		Description:  "main node",
		PluginName:   "node",
		AutoLoad:     true,
		RestoreState: true,
		Attributes:   node.NewAttr(),
		Settings:     settings,
	}
}

// GetNewMoon ...
func GetNewMoon(name string) *m.Entity {
	settings := moon.NewSettings()
	settings[moon.AttrLat].Value = 54.9022
	settings[moon.AttrLon].Value = 83.0335
	return &m.Entity{
		Id:           common.EntityId(fmt.Sprintf("moon.%s", name)),
		Description:  "home",
		PluginName:   "moon",
		AutoLoad:     true,
		RestoreState: true,
		Attributes:   moon.NewAttr(),
		Settings:     settings,
		States: []*m.EntityState{
			{
				Name:        moon.StateAboveHorizon,
				Description: "above horizon",
			},
			{
				Name:        moon.StateBelowHorizon,
				Description: "below horizon",
			},
		},
	}
}

// GetNewWeatherMet ...
func GetNewWeatherMet(name string) *m.Entity {
	settings := weather.NewSettings()
	settings[weather.AttrLat].Value = 54.9022
	settings[weather.AttrLon].Value = 83.0335
	return &m.Entity{
		Id:           common.EntityId(fmt.Sprintf("weather_met.%s", name)),
		Description:  name,
		PluginName:   "weather_met",
		AutoLoad:     true,
		RestoreState: true,
		Attributes:   weather.BaseForecast(),
		Settings:     settings,
	}
}

// GetNewWeatherOwm ...
//func GetNewWeatherOwm(name string) *m.Entity {
//	settings := weather_owm.NewSettings()
//	settings[weather_owm.AttrAppid].Value = "**************"
//	settings[weather_owm.AttrUnits].Value = "metric"
//	settings[weather_owm.AttrLang].Value = "ru"
//	return &m.Entity{
//		Id:          common.EntityId(fmt.Sprintf("weather_owm.%s", name)),
//		Description: "weather owm",
//		PluginName:  weather_owm.EntityWeatherOwm,
//		AutoLoad:    true,
//		Settings:    settings,
//	}
//}

// GetNewSun ...
func GetNewSun(name string) *m.Entity {
	settings := sun.NewSettings()
	settings[sun.AttrLat].Value = 54.9022
	settings[sun.AttrLon].Value = 83.0335
	return &m.Entity{
		Id:           common.EntityId(fmt.Sprintf("sun.%s", name)),
		Description:  "home",
		PluginName:   "sun",
		AutoLoad:     true,
		RestoreState: true,
		Attributes:   sun.NewAttr(),
		Settings:     settings,
		States: []*m.EntityState{
			{
				Name:        sun.AttrDusk,
				Description: "dusk (evening nautical twilight starts)",
			},
		},
	}
}

// GetNewBitmineL3 ...
func GetNewBitmineL3(name string) *m.Entity {
	settings := cgminer.NewSettings()
	settings[cgminer.SettingHost].Value = "192.168.0.243"
	settings[cgminer.SettingPort].Value = 4028
	settings[cgminer.SettingTimeout].Value = 2
	settings[cgminer.SettingUser].Value = "user"
	settings[cgminer.SettingPass].Value = "pass"
	settings[cgminer.SettingManufacturer].Value = bitmine.ManufactureBitmine
	settings[cgminer.SettingModel].Value = bitmine.DeviceL3Plus
	return &m.Entity{
		Id:           common.EntityId(fmt.Sprintf("cgminer.%s", name)),
		Description:  "antminer L3",
		PluginName:   "cgminer",
		AutoLoad:     true,
		RestoreState: true,
		Attributes:   cgminer.NewAttr(),
		Settings:     settings,
	}
}

// GetNewSensor ...
func GetNewSensor(name string) *m.Entity {

	return &m.Entity{
		Id:           common.EntityId(fmt.Sprintf("sensor.%s", name)),
		Description:  "api",
		PluginName:   "sensor",
		AutoLoad:     true,
		RestoreState: true,
	}
}

// GetNewModbusRtu ...
func GetNewModbusRtu(name string) *m.Entity {
	return &m.Entity{
		Id:           common.EntityId(fmt.Sprintf("modbus_rtu.%s", name)),
		Description:  fmt.Sprintf("%s entity", name),
		PluginName:   "modbus_rtu",
		AutoLoad:     true,
		RestoreState: true,
		Attributes:   modbus_rtu.NewAttr(),
		Settings:     modbus_rtu.NewSettings(),
	}
}

// GetNewModbusTcp ...
func GetNewModbusTcp(name string) *m.Entity {
	return &m.Entity{
		Id:           common.EntityId(fmt.Sprintf("modbus_tcp.%s", name)),
		Description:  fmt.Sprintf("%s entity", name),
		PluginName:   "modbus_tcp",
		AutoLoad:     true,
		RestoreState: true,
		Attributes:   modbus_tcp.NewAttr(),
		Settings:     modbus_tcp.NewSettings(),
	}
}

// GetNewTelegram ...
func GetNewTelegram(name string) *m.Entity {
	settings := telegram.NewSettings()
	settings[telegram.AttrToken].Value = "XXXX"
	return &m.Entity{
		Id:           common.EntityId(fmt.Sprintf("%s.%s", telegram.Name, name)),
		Description:  "",
		PluginName:   telegram.Name,
		AutoLoad:     true,
		RestoreState: true,
		Attributes:   telegram.NewAttr(),
		Settings:     settings,
	}
}

// AddPlugin ...
func AddPlugin(adaptors *adaptors.Adaptors, name string, opts ...m.AttributeValue) (err error) {
	plugin := &m.Plugin{
		Name:    name,
		Version: "0.0.1",
		Enabled: true,
		System:  true,
	}
	if len(opts) > 0 {
		plugin.Settings = opts[0]
	}
	err = adaptors.Plugin.CreateOrUpdate(context.Background(), plugin)
	return
}

// RegisterConvey ...
func RegisterConvey(scriptService scripts.ScriptService, ctx convey.C) {
	scriptService.PushFunctions("So", func(actual interface{}, assert string, expected interface{}) {
		//fmt.Printf("actual(%v), expected(%v)\n", actual, expected)
		switch assert {
		case "ShouldEqual":
			ctx.So(fmt.Sprintf("%v", actual), convey.ShouldEqual, expected)
		case "ShouldNotBeBlank":
			ctx.So(fmt.Sprintf("%v", actual), convey.ShouldNotBeBlank)
		}
	})

}

// Wait ...
func Wait(timeOut time.Duration, ch chan interface{}) (ok bool) {

	select {
	case <-ch:
		ok = true
	case <-time.After(timeOut):
	}
	return
}

// WaitT ...
func WaitT[T events.EventTriggerLoaded | events.EventTriggerCompleted | events.EventStateChanged | alexa.EventAlexaAction | []byte | struct{}](timeOut time.Duration, ch chan T) (v T, ok bool) {

	select {
	case v = <-ch:
		ok = true
		return
	case <-time.After(timeOut):
	}
	return
}

type accepted struct {
	conn net.Conn
	err  error
}

// MockHttpServer ...
func MockHttpServer(ctx context.Context, ip string, port int64, payload []byte) (err error) {

	var listener net.Listener
	if listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port)); err != nil {
		return
	}

	_ = http.Serve(listener, http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		_, _ = fmt.Fprintf(rw, string(payload))
	}))
	c := make(chan accepted, 1)
	for {
		select {
		case <-ctx.Done():
			_ = listener.Close()
			return
		case a := <-c:
			if a.err != nil {
				err = a.err
				return
			}
			go func(conn net.Conn) {
				_, _ = conn.Write(payload)
				_ = conn.Close()
			}(a.conn)
		default:
		}
	}
}

// MockTCPServer ...
func MockTCPServer(ctx context.Context, ip string, port int64, payloads ...[]byte) (err error) {
	var listener net.Listener
	if listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port)); err != nil {
		return
	}
	c := make(chan accepted, 3)
	go func() {
		for {
			conn, err := listener.Accept()
			c <- accepted{conn, err}
		}
	}()
	var counter int
	for {
		select {
		case <-ctx.Done():
			_ = listener.Close()
			return
		case a := <-c:
			if a.err != nil {
				err = a.err
				return
			}
			go func(conn net.Conn) {
				if counter < len(payloads) {
					_, _ = conn.Write(payloads[counter])
				} else {
					_, _ = conn.Write(payloads[len(payloads)-1])
				}
				_ = conn.Close()
				counter++
			}(a.conn)
		default:
		}
	}
}

// GetPort ...
func GetPort() int64 {
	port, _ := freeport.GetFreePort()
	return int64(port)
}

// AddScript ...
func AddScript(name, src string, adaptors *adaptors.Adaptors, scriptService scripts.ScriptService) (script *m.Script, err error) {

	script = &m.Script{
		Lang:        common.ScriptLangCoffee,
		Name:        name,
		Source:      src,
		Description: "description " + name,
	}

	var engine *scripts.Engine
	if engine, err = scriptService.NewEngine(script); err != nil {
		return
	}

	if err = engine.Compile(); err != nil {
		return
	}

	script.Id, err = adaptors.Script.Add(context.Background(), script)

	return
}

func AddTrigger(trigger *m.NewTrigger, adaptors *adaptors.Adaptors, eventBus bus.Bus) (id int64, err error) {
	if id, err = adaptors.Trigger.Add(context.Background(), trigger); err != nil {
		return
	}
	eventBus.Publish(fmt.Sprintf("system/automation/triggers/%d", id), events.EventCreatedTriggerModel{
		Id: id,
	})
	return
}

func AddTask(newTask *m.NewTask, adaptors *adaptors.Adaptors, eventBus bus.Bus) (task1Id int64, err error) {
	if task1Id, err = adaptors.Task.Add(context.Background(), newTask); err != nil {
		return
	}
	eventBus.Publish(fmt.Sprintf("system/automation/tasks/%d", task1Id), events.EventCreatedTaskModel{
		Id: task1Id,
	})
	return
}

func WaitTask(eventBus bus.Bus, timeOut time.Duration, tasks ...int64) (result chan bool) {

	list := map[int64]bool{}
	for _, task := range tasks {
		list[task] = false
	}

	var closed = atomic.NewBool(false)
	result = make(chan bool, 1)
	go func() {
		mx := sync.Mutex{}

		ch := make(chan interface{})
		defer close(ch)
		fn := func(_ string, msg interface{}) {
			switch v := msg.(type) {
			case events.EventTaskLoaded:
				mx.Lock()
				defer mx.Unlock()
				fmt.Printf("Task %d loaded ...\r\n", v.Id)
				if _, ok := list[v.Id]; ok {
					list[v.Id] = true
				} else {
					return
				}
				for _, loaded := range list {
					if !loaded {
						return
					}
				}
				if closed.Load() {
					return
				}
				ch <- struct{}{}
			}

		}
		eventBus.Subscribe("system/automation/tasks/+", fn, true)
		defer eventBus.Unsubscribe("system/automation/tasks/+", fn)

		result <- Wait(timeOut, ch)
		closed.Store(true)
		close(result)

	}()

	return
}

func WaitEntity(eventBus bus.Bus, timeOut time.Duration, entities ...string) (result chan bool) {

	list := map[string]bool{}
	for _, entity := range entities {
		list[entity] = false
	}

	var closed = atomic.NewBool(false)
	result = make(chan bool, 1)
	go func() {
		mx := sync.Mutex{}

		ch := make(chan interface{})
		defer close(ch)
		fn := func(_ string, msg interface{}) {
			switch v := msg.(type) {
			case events.EventEntityLoaded:
				mx.Lock()
				defer mx.Unlock()
				fmt.Printf("Plugin %s loaded ...\r\n", v.EntityId.String())
				if _, ok := list[v.EntityId.String()]; ok {
					list[v.EntityId.String()] = true
				} else {
					return
				}
				for _, loaded := range list {
					if !loaded {
						return
					}
				}
				if closed.Load() {
					return
				}
				ch <- struct{}{}
			}

		}
		eventBus.Subscribe("system/entities/+", fn, true)
		defer eventBus.Unsubscribe("system/entities/+", fn)

		result <- Wait(timeOut, ch)
		closed.Store(true)
		close(result)

	}()

	return
}

func WaitPlugins(eventBus bus.Bus, timeOut time.Duration, plugins ...string) (result chan bool) {

	list := map[string]bool{}
	for _, plugin := range plugins {
		list[plugin] = false
	}

	var closed = atomic.NewBool(false)
	result = make(chan bool, 1)
	go func() {
		mx := sync.Mutex{}

		ch := make(chan interface{})
		defer close(ch)
		fn := func(_ string, msg interface{}) {
			switch v := msg.(type) {
			case events.EventPluginLoaded:
				mx.Lock()
				defer mx.Unlock()
				fmt.Printf("Plugin %s loaded ...\r\n", v.PluginName)
				if _, ok := list[v.PluginName]; ok {
					list[v.PluginName] = true
				} else {
					return
				}
				for _, loaded := range list {
					if !loaded {
						return
					}
				}
				if closed.Load() {
					return
				}
				ch <- struct{}{}
			}

		}
		eventBus.Subscribe("system/plugins/+", fn, true)
		defer eventBus.Unsubscribe("system/plugins/+", fn)

		result <- Wait(timeOut, ch)
		closed.Store(true)
		close(result)

	}()

	return
}

func WaitService(eventBus bus.Bus, timeOut time.Duration, services ...string) (result chan bool) {

	list := map[string]bool{}
	for _, service := range services {
		list[service] = false
	}

	var closed = atomic.NewBool(false)
	result = make(chan bool, 1)
	go func() {
		mx := sync.Mutex{}

		ch := make(chan interface{})
		defer close(ch)
		fn := func(_ string, msg interface{}) {
			switch v := msg.(type) {
			case events.EventServiceStarted:
				mx.Lock()
				defer mx.Unlock()
				fmt.Printf("Service %s started ...\r\n", v.Service)
				if _, ok := list[v.Service]; ok {
					list[v.Service] = true
				} else {
					return
				}
				for _, started := range list {
					if !started {
						return
					}
				}
				if closed.Load() {
					return
				}
				ch <- struct{}{}
			}

		}
		eventBus.Subscribe("system/services/+", fn, true)
		defer eventBus.Unsubscribe("system/services/+", fn)

		time.Sleep(time.Millisecond * 500)

		result <- Wait(timeOut, ch)
		closed.Store(true)
		close(result)

	}()

	return
}

func WaitMessage[T events.EventStateChanged | events.EventTriggerCompleted | events.EventTriggerLoaded](
	eventBus bus.Bus, timeOut time.Duration, topic string, options ...interface{},
) (msg T, ok bool) {

	var closed = atomic.NewBool(false)
	ch := make(chan T)
	defer close(ch)
	defer closed.Store(true)
	fn := func(_ string, msg interface{}) {
		switch v := msg.(type) {
		case T:
			if closed.Load() {
				return
			}
			ch <- v
		}
	}

	var retain bool
	if len(options) > 0 {
		retain, _ = options[0].(bool)
	}

	eventBus.Subscribe(topic, fn, retain)
	defer eventBus.Unsubscribe(topic, fn)

	msg, ok = WaitT[T](timeOut, ch)

	time.Sleep(time.Millisecond * 500)

	return
}

func WaitGroupTimeout(wg *sync.WaitGroup, timeOut time.Duration) bool {
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		wg.Wait()
	}()
	select {
	case <-ch:
		return true
	case <-time.After(timeOut):
		return false
	}
}
