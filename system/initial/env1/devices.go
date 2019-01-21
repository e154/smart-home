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

	// action1
	deviceAction1 := &m.DeviceAction{
		Name:     "Condition check",
		DeviceId: device1.Id,
		ScriptId: scripts["mb_dev1_condition_check_v1"].Id,
	}
	ok, _ = deviceAction1.Valid()
	So(ok, ShouldEqual, true)
	deviceAction1.Id, err = adaptors.DeviceAction.Add(deviceAction1)
	So(err, ShouldBeNil)
	deviceActions = append(deviceActions, deviceAction1)

	// action2
	deviceAction2 := &m.DeviceAction{
		Name:     "turn on light1",
		DeviceId: device1.Id,
		ScriptId: scripts["mb_dev1_turn_on_first_light_v1"].Id,
	}
	ok, _ = deviceAction2.Valid()
	So(ok, ShouldEqual, true)
	deviceAction2.Id, err = adaptors.DeviceAction.Add(deviceAction2)
	So(err, ShouldBeNil)
	deviceActions = append(deviceActions, deviceAction2)

	// action3
	deviceAction3 := &m.DeviceAction{
		Name:     "turn off light1",
		DeviceId: device1.Id,
		ScriptId: scripts["mb_dev1_turn_off_first_light_v1"].Id,
	}
	ok, _ = deviceAction3.Valid()
	So(ok, ShouldEqual, true)
	deviceAction3.Id, err = adaptors.DeviceAction.Add(deviceAction3)
	So(err, ShouldBeNil)
	deviceActions = append(deviceActions, deviceAction3)

	// states
	stateDev1Enabled := &m.DeviceState{
		SystemName:  "ENABLED",
		Description: "device enabled",
		DeviceId:    device1.Id,
	}
	stateDev1Disabled := &m.DeviceState{
		SystemName:  "DISABLED",
		Description: "device disabled",
		DeviceId:    device1.Id,
	}
	stateDev1Error := &m.DeviceState{
		SystemName:  "ERROR",
		Description: "device in error state",
		DeviceId:    device1.Id,
	}
	stateDev1Light1On := &m.DeviceState{
		SystemName:  "LIGHT_1_ON",
		Description: "device light 1 on",
		DeviceId:    device1.Id,
	}
	stateDev1Light1Off := &m.DeviceState{
		SystemName:  "LIGHT_1_OFF",
		Description: "device light 1 off",
		DeviceId:    device1.Id,
	}
	stateDev1Light2On := &m.DeviceState{
		SystemName:  "LIGHT_2_ON",
		Description: "device light 2 on",
		DeviceId:    device1.Id,
	}
	stateDev1Light2Off := &m.DeviceState{
		SystemName:  "LIGHT_2_OFF",
		Description: "device light 2 off",
		DeviceId:    device1.Id,
	}
	stateDev1Light3On := &m.DeviceState{
		SystemName:  "LIGHT_3_ON",
		Description: "device light 3 on",
		DeviceId:    device1.Id,
	}
	stateDev1Light3Off := &m.DeviceState{
		SystemName:  "LIGHT_3_OFF",
		Description: "device light 3 off",
		DeviceId:    device1.Id,
	}
	stateDev1Light4On := &m.DeviceState{
		SystemName:  "LIGHT_4_ON",
		Description: "device light 4 on",
		DeviceId:    device1.Id,
	}
	stateDev1Light4Off := &m.DeviceState{
		SystemName:  "LIGHT_4_OFF",
		Description: "device light 4 off",
		DeviceId:    device1.Id,
	}
	stateDev1Fan1On := &m.DeviceState{
		SystemName:  "FAN_1_ON",
		Description: "device fan 1 on",
		DeviceId:    device1.Id,
	}
	stateDev1Fan1Off := &m.DeviceState{
		SystemName:  "FAN_1_OFF",
		Description: "device fan 1 off",
		DeviceId:    device1.Id,
	}
	ok, _ = stateDev1Enabled.Valid()
	So(ok, ShouldEqual, true)
	ok, _ = stateDev1Disabled.Valid()
	So(ok, ShouldEqual, true)
	ok, _ = stateDev1Error.Valid()
	So(ok, ShouldEqual, true)
	ok, _ = stateDev1Light2On.Valid()
	So(ok, ShouldEqual, true)
	ok, _ = stateDev1Light2Off.Valid()
	So(ok, ShouldEqual, true)
	ok, _ = stateDev1Light3On.Valid()
	So(ok, ShouldEqual, true)
	ok, _ = stateDev1Light3Off.Valid()
	So(ok, ShouldEqual, true)
	ok, _ = stateDev1Light4On.Valid()
	So(ok, ShouldEqual, true)
	ok, _ = stateDev1Light4Off.Valid()
	So(ok, ShouldEqual, true)
	ok, _ = stateDev1Fan1On.Valid()
	So(ok, ShouldEqual, true)
	ok, _ = stateDev1Fan1Off.Valid()
	So(ok, ShouldEqual, true)

	stateDev1Enabled.Id, err = adaptors.DeviceState.Add(stateDev1Enabled)
	So(err, ShouldBeNil)
	stateDev1Disabled.Id, err = adaptors.DeviceState.Add(stateDev1Disabled)
	So(err, ShouldBeNil)
	stateDev1Error.Id, err = adaptors.DeviceState.Add(stateDev1Error)
	So(err, ShouldBeNil)
	stateDev1Light1On.Id, err = adaptors.DeviceState.Add(stateDev1Light1On)
	So(err, ShouldBeNil)
	stateDev1Light1Off.Id, err = adaptors.DeviceState.Add(stateDev1Light1Off)
	So(err, ShouldBeNil)
	stateDev1Light2On.Id, err = adaptors.DeviceState.Add(stateDev1Light2On)
	So(err, ShouldBeNil)
	stateDev1Light2Off.Id, err = adaptors.DeviceState.Add(stateDev1Light2Off)
	So(err, ShouldBeNil)
	stateDev1Light3On.Id, err = adaptors.DeviceState.Add(stateDev1Light3On)
	So(err, ShouldBeNil)
	stateDev1Light3Off.Id, err = adaptors.DeviceState.Add(stateDev1Light3Off)
	So(err, ShouldBeNil)
	stateDev1Light4On.Id, err = adaptors.DeviceState.Add(stateDev1Light4On)
	So(err, ShouldBeNil)
	stateDev1Light4Off.Id, err = adaptors.DeviceState.Add(stateDev1Light4Off)
	So(err, ShouldBeNil)
	stateDev1Fan1On.Id, err = adaptors.DeviceState.Add(stateDev1Fan1On)
	So(err, ShouldBeNil)
	stateDev1Fan1Off.Id, err = adaptors.DeviceState.Add(stateDev1Fan1Off)
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

	deviceAction21 := &m.DeviceAction{
		Name:     "internet address condition check",
		DeviceId: device2.Id,
		ScriptId: scripts["cmd_condition_check_v1"].Id,
	}
	ok, _ = deviceAction21.Valid()
	So(ok, ShouldEqual, true)
	deviceAction21.Id, err = adaptors.DeviceAction.Add(deviceAction21)
	So(err, ShouldBeNil)

	deviceActions = append(deviceActions, deviceAction21)

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
