package env1

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/system/initial/assertions"
	"github.com/e154/smart-home/common"
)

func devices(node1 *m.Node,
	adaptors *adaptors.Adaptors,
	script1, script2, script3, script7 *m.Script) (device1, device2, device3 *m.Device, deviceAction3, deviceAction7 *m.DeviceAction) {

	// devices
	// ------------------------------------------------
	device1 = &m.Device{
		Name:       "device1",
		Status:     "enabled",
		Type:       common.DevTypeDefault,
		Node:       node1,
		IsGroup:    true,
		Properties: []byte("{}"),
	}

	smartBusConfig := &common.DevConfSmartBus{
		Baud:     19200,
		Device:   0,
		Timeout:  457,
		StopBits: 2,
		Sleep:    0,
	}

	ok, _ := device1.SetProperties(smartBusConfig)
	So(ok, ShouldEqual, true)

	ok, _ = device1.Valid()
	So(ok, ShouldEqual, true)

	device3 = &m.Device{
		Name:       "device3",
		Status:     "disabled",
		Type:       common.DevTypeDefault,
		Node:       node1,
		Properties: []byte("{}"),
	}
	ok, _ = device3.Valid()
	So(ok, ShouldEqual, true)

	var err error
	device1.Id, err = adaptors.Device.Add(device1)
	So(err, ShouldBeNil)
	device3.Id, err = adaptors.Device.Add(device3)
	So(err, ShouldBeNil)

	device2 = &m.Device{
		Name:       "device2",
		Status:     "enabled",
		Type:       common.DevTypeDefault,
		Device:     device1,
		Properties: []byte("{}"),
	}

	ok, _ = device2.Valid()
	So(ok, ShouldEqual, true)

	smartBusConfig2 := &common.DevConfSmartBus{
		Baud:     19200,
		Device:   2,
		Timeout:  457,
		StopBits: 2,
		Sleep:    0,
	}

	ok, _ = device2.SetProperties(smartBusConfig2)
	So(ok, ShouldEqual, true)

	device2.Id, err = adaptors.Device.Add(device2)
	So(err, ShouldBeNil)

	device4 := &m.Device{
		Name:       "device4",
		Status:     "enabled",
		Type:       common.DevTypeDefault,
		Device:     device1,
		Properties: []byte("{}"),
	}

	ok, _ = device4.Valid()
	So(ok, ShouldEqual, true)

	smartBusConfig4 := &common.DevConfSmartBus{
		Baud:     19200,
		Device:   1,
		Timeout:  457,
		StopBits: 2,
		Sleep:    0,
	}

	ok, _ = device4.SetProperties(smartBusConfig4)
	So(ok, ShouldEqual, true)

	device4.Id, err = adaptors.Device.Add(device4)
	So(err, ShouldBeNil)

	// add device action
	// ------------------------------------------------
	deviceAction1 := &m.DeviceAction{
		Name:     "Turning on",
		DeviceId: device1.Id,
		ScriptId: script1.Id,
	}
	deviceAction2 := &m.DeviceAction{
		Name:     "Power off",
		DeviceId: device1.Id,
		ScriptId: script2.Id,
	}
	deviceAction3 = &m.DeviceAction{
		Name:     "Condition check",
		DeviceId: device1.Id,
		ScriptId: script3.Id,
	}
	ok, _ = deviceAction1.Valid()
	So(ok, ShouldEqual, true)
	ok, _ = deviceAction2.Valid()
	So(ok, ShouldEqual, true)
	ok, _ = deviceAction3.Valid()
	So(ok, ShouldEqual, true)

	deviceAction1.Id, err = adaptors.DeviceAction.Add(deviceAction1)
	So(err, ShouldBeNil)
	deviceAction2.Id, err = adaptors.DeviceAction.Add(deviceAction2)
	So(err, ShouldBeNil)
	deviceAction3.Id, err = adaptors.DeviceAction.Add(deviceAction3)
	So(err, ShouldBeNil)

	deviceAction4 := &m.DeviceAction{
		Name:     "Turning on",
		DeviceId: device1.Id,
		ScriptId: script1.Id,
	}
	deviceAction5 := &m.DeviceAction{
		Name:     "Power off",
		DeviceId: device1.Id,
		ScriptId: script2.Id,
	}
	deviceAction6 := &m.DeviceAction{
		Name:     "Condition check",
		DeviceId: device1.Id,
		ScriptId: script3.Id,
	}

	ok, _ = deviceAction4.Valid()
	So(ok, ShouldEqual, true)
	ok, _ = deviceAction5.Valid()
	So(ok, ShouldEqual, true)
	ok, _ = deviceAction6.Valid()
	So(ok, ShouldEqual, true)

	deviceAction4.Id, err = adaptors.DeviceAction.Add(deviceAction4)
	So(err, ShouldBeNil)
	deviceAction5.Id, err = adaptors.DeviceAction.Add(deviceAction5)
	So(err, ShouldBeNil)
	deviceAction6.Id, err = adaptors.DeviceAction.Add(deviceAction6)
	So(err, ShouldBeNil)

	// add device state
	// ------------------------------------------------

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

	deviceState4 := &m.DeviceState{
		SystemName:  "ENABLED",
		Description: "device enabled",
		DeviceId:    device3.Id,
	}
	deviceState5 := &m.DeviceState{
		SystemName:  "DISABLED",
		Description: "device disabled",
		DeviceId:    device3.Id,
	}
	deviceState6 := &m.DeviceState{
		SystemName:  "ERROR",
		Description: "device in error state",
		DeviceId:    device3.Id,
	}
	ok, _ = deviceState4.Valid()
	So(ok, ShouldEqual, true)
	ok, _ = deviceState5.Valid()
	So(ok, ShouldEqual, true)
	ok, _ = deviceState6.Valid()
	So(ok, ShouldEqual, true)

	deviceState4.Id, err = adaptors.DeviceState.Add(deviceState4)
	So(err, ShouldBeNil)
	deviceState5.Id, err = adaptors.DeviceState.Add(deviceState5)
	So(err, ShouldBeNil)
	deviceState6.Id, err = adaptors.DeviceState.Add(deviceState6)
	So(err, ShouldBeNil)

	// device for command action
	// ------------------------------------------------

	device6 := &m.Device{
		Name:       "device6",
		Status:     "enabled",
		Type:       common.DevTypeCommand,
		Node:       node1,
		Properties: []byte("{}"),
	}

	ok, _ = device6.Valid()
	So(ok, ShouldEqual, true)
	device6.Id, err = adaptors.Device.Add(device6)
	So(err, ShouldBeNil)

	deviceAction7 = &m.DeviceAction{
		Name:     "Condition check",
		DeviceId: device6.Id,
		ScriptId: script7.Id,
	}
	ok, _ = deviceAction7.Valid()
	So(ok, ShouldEqual, true)
	deviceAction7.Id, err = adaptors.DeviceAction.Add(deviceAction7)
	So(err, ShouldBeNil)

	deviceState7 := &m.DeviceState{
		SystemName:  "ONLINE",
		Description: "address is online",
		DeviceId:    device6.Id,
	}
	deviceState8 := &m.DeviceState{
		SystemName:  "OFFLINE",
		Description: "address is offline",
		DeviceId:    device6.Id,
	}
	deviceState9 := &m.DeviceState{
		SystemName:  "ERROR",
		Description: "unknown error",
		DeviceId:    device6.Id,
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

	return
}
