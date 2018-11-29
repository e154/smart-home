package env1

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/system/initial/assertions"
)

func devices(node1 *m.Node,
	adaptors *adaptors.Adaptors,
	script1, script2, script3 *m.Script) (device1, device2, device3 *m.Device, deviceAction3 *m.DeviceAction) {

	// devices
	// ------------------------------------------------
	device1 = &m.Device{
		Name:       "device1",
		Status:     "enabled",
		Type:       "default",
		Node:       node1,
		Properties: []byte("{}"),
	}

	ok, _ := device1.Valid()
	So(ok, ShouldEqual, true)

	device3 = &m.Device{
		Name:       "device3",
		Status:     "disabled",
		Type:       "default",
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
		Type:       "default",
		Device:     device1,
		Properties: []byte("{}"),
	}

	ok, _ = device2.Valid()
	So(ok, ShouldEqual, true)
	device2.Id, err = adaptors.Device.Add(device2)
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

	return
}
