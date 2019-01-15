package env1

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/system/initial/assertions"
	. "github.com/e154/smart-home/models/devices"
)

func devices(node1 *m.Node,
	adaptors *adaptors.Adaptors,
	scripts map[string]*m.Script) (devices []*m.Device, deviceActions []*m.DeviceAction) {

	// devices1
	// ------------------------------------------------
	device1 := &m.Device{
		Name:       "device1",
		Status:     "enabled",
		Type:       DevTypeModbus,
		Node:       node1,
		Properties: []byte("{}"),
	}

	modBusConfig := &DevModBusConfig{
		SlaveId:  1,
		Baud:     19200,
		DataBits: 8,
		StopBits: 1,
		Parity:   "none",
		Timeout:  100,
	}

	ok, _ := device1.SetProperties(modBusConfig)
	So(ok, ShouldEqual, true)

	ok, _ = device1.Valid()
	So(ok, ShouldEqual, true)

	var err error
	device1.Id, err = adaptors.Device.Add(device1)
	So(err, ShouldBeNil)

	devices = append(devices, device1)

	deviceAction1 := &m.DeviceAction{
		Name:     "Condition check",
		DeviceId: device1.Id,
		ScriptId: scripts["mb_condition_check_v1"].Id,
	}

	ok, _ = deviceAction1.Valid()
	So(ok, ShouldEqual, true)

	deviceAction1.Id, err = adaptors.DeviceAction.Add(deviceAction1)
	So(err, ShouldBeNil)

	deviceActions = append(deviceActions, deviceAction1)

	deviceState1 := &m.DeviceState{
		SystemName:  "ENABLED",
		Description: "device enabled",
		DeviceId:    device1.Id,
	}
	deviceState2 := &m.DeviceState{
		SystemName:  "DISABLED",
		Description: "device disabled",
		DeviceId:    device1.Id,
	}
	deviceState3 := &m.DeviceState{
		SystemName:  "ERROR",
		Description: "device in error state",
		DeviceId:    device1.Id,
	}
	ok, _ = deviceState1.Valid()
	So(ok, ShouldEqual, true)
	ok, _ = deviceState2.Valid()
	So(ok, ShouldEqual, true)
	ok, _ = deviceState3.Valid()
	So(ok, ShouldEqual, true)

	deviceState1.Id, err = adaptors.DeviceState.Add(deviceState1)
	So(err, ShouldBeNil)
	deviceState2.Id, err = adaptors.DeviceState.Add(deviceState2)
	So(err, ShouldBeNil)
	deviceState3.Id, err = adaptors.DeviceState.Add(deviceState3)
	So(err, ShouldBeNil)

	// device 2
	// ------------------------------------------------
	device2 := &m.Device{
		Name:       "device2",
		Status:     "enabled",
		Type:       DevTypeCommand,
		Node:       node1,
		Properties: []byte("{}"),
	}

	ok, _ = device2.Valid()
	So(ok, ShouldEqual, true)
	device2.Id, err = adaptors.Device.Add(device2)
	So(err, ShouldBeNil)

	devices = append(devices, device2)

	deviceAction2 := &m.DeviceAction{
		Name:     "internet address condition check",
		DeviceId: device2.Id,
		ScriptId: scripts["cmd_condition_check_v1"].Id,
	}
	ok, _ = deviceAction2.Valid()
	So(ok, ShouldEqual, true)
	deviceAction2.Id, err = adaptors.DeviceAction.Add(deviceAction2)
	So(err, ShouldBeNil)

	deviceActions = append(deviceActions, deviceAction2)

	deviceState7 := &m.DeviceState{
		SystemName:  "ONLINE",
		Description: "address is online",
		DeviceId:    device2.Id,
	}
	deviceState8 := &m.DeviceState{
		SystemName:  "OFFLINE",
		Description: "address is offline",
		DeviceId:    device2.Id,
	}
	deviceState9 := &m.DeviceState{
		SystemName:  "ERROR",
		Description: "unknown error",
		DeviceId:    device2.Id,
	}
	ok, _ = deviceState7.Valid()
	So(ok, ShouldEqual, true)
	ok, _ = deviceState8.Valid()
	So(ok, ShouldEqual, true)
	ok, _ = deviceState9.Valid()
	So(ok, ShouldEqual, true)

	deviceState7.Id, err = adaptors.DeviceState.Add(deviceState7)
	So(err, ShouldBeNil)
	deviceState8.Id, err = adaptors.DeviceState.Add(deviceState8)
	So(err, ShouldBeNil)
	deviceState9.Id, err = adaptors.DeviceState.Add(deviceState9)
	So(err, ShouldBeNil)

	// devices3
	// ------------------------------------------------
	device3 := &m.Device{
		Name:       "device3",
		Status:     "enabled",
		Type:       DevTypeModbus,
		Node:       node1,
		IsGroup:    true,
		Properties: []byte("{}"),
	}

	modBusConfig = &DevModBusConfig{
		Baud:     115200,
		DataBits: 8,
		StopBits: 2,
		Parity:   "none",
	}

	ok, _ = device3.SetProperties(modBusConfig)
	So(ok, ShouldEqual, true)

	ok, _ = device3.Valid()
	So(ok, ShouldEqual, true)

	device3.Id, err = adaptors.Device.Add(device3)
	So(err, ShouldBeNil)

	devices = append(devices, device3)

	// devices4
	// ------------------------------------------------
	device4 := &m.Device{
		Name:       "device4",
		Status:     "enabled",
		Type:       DevTypeModbus,
		Device:     device3,
		Properties: []byte("{\"slave_id\": 2}"),
	}

	ok, _ = device4.Valid()
	So(ok, ShouldEqual, true)

	device4.Id, err = adaptors.Device.Add(device4)
	So(err, ShouldBeNil)

	devices = append(devices, device4)

	return
}
