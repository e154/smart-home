package env1

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/models/devices"
	. "github.com/e154/smart-home/system/initial/assertions"
)

func devices(node1 *m.Node,
	adaptors *adaptors.Adaptors,
	scripts map[string]*m.Script) (devices []*m.Device, deviceActions map[string]*m.DeviceAction, deviceStates map[string]*m.DeviceState) {

	deviceActions = make(map[string]*m.DeviceAction)
	deviceStates = make(map[string]*m.DeviceState)

	// devices1
	// ------------------------------------------------
	device1 := &m.Device{
		Name:       "device1",
		Status:     "enabled",
		Type:       DevTypeModbusTcp,
		Node:       node1,
		Properties: []byte("{}"),
	}

	modBusConfig := &DevModBusTcpConfig{
		//SlaveId:  1,
		//Baud:     19200,
		//DataBits: 8,
		//StopBits: 1,
		//Parity:   "none",
		//Timeout:  100,
		AddressPort: "127.0.0.1:502",
		SlaveId: 1,
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
	deviceActions["mb_dev1_condition_check_v1"] = deviceAction1

	// light1
	deviceAction2 := &m.DeviceAction{
		Name:     "turn on light1",
		DeviceId: device1.Id,
		ScriptId: scripts["mb_dev1_turn_on_light1_v1"].Id,
	}
	ok, _ = deviceAction2.Valid()
	So(ok, ShouldEqual, true)
	deviceAction2.Id, err = adaptors.DeviceAction.Add(deviceAction2)
	So(err, ShouldBeNil)
	deviceActions["mb_dev1_turn_on_light1_v1"] = deviceAction2

	deviceAction3 := &m.DeviceAction{
		Name:     "turn off light1",
		DeviceId: device1.Id,
		ScriptId: scripts["mb_dev1_turn_off_light1_v1"].Id,
	}
	ok, _ = deviceAction3.Valid()
	So(ok, ShouldEqual, true)
	deviceAction3.Id, err = adaptors.DeviceAction.Add(deviceAction3)
	So(err, ShouldBeNil)
	deviceActions["mb_dev1_turn_off_light1_v1"] = deviceAction3

	// light2
	deviceAction4 := &m.DeviceAction{
		Name:     "turn on light2",
		DeviceId: device1.Id,
		ScriptId: scripts["mb_dev1_turn_on_light2_v1"].Id,
	}
	ok, _ = deviceAction4.Valid()
	So(ok, ShouldEqual, true)
	deviceAction4.Id, err = adaptors.DeviceAction.Add(deviceAction4)
	So(err, ShouldBeNil)
	deviceActions["mb_dev1_turn_on_light2_v1"] = deviceAction4

	deviceAction5 := &m.DeviceAction{
		Name:     "turn off light2",
		DeviceId: device1.Id,
		ScriptId: scripts["mb_dev1_turn_off_light2_v1"].Id,
	}
	ok, _ = deviceAction5.Valid()
	So(ok, ShouldEqual, true)
	deviceAction5.Id, err = adaptors.DeviceAction.Add(deviceAction5)
	So(err, ShouldBeNil)
	deviceActions["mb_dev1_turn_off_light2_v1"] = deviceAction5

	// light3
	deviceAction6 := &m.DeviceAction{
		Name:     "turn on light3",
		DeviceId: device1.Id,
		ScriptId: scripts["mb_dev1_turn_on_light3_v1"].Id,
	}
	ok, _ = deviceAction6.Valid()
	So(ok, ShouldEqual, true)
	deviceAction6.Id, err = adaptors.DeviceAction.Add(deviceAction6)
	So(err, ShouldBeNil)
	deviceActions["mb_dev1_turn_on_light3_v1"] = deviceAction6

	deviceAction7 := &m.DeviceAction{
		Name:     "turn off light3",
		DeviceId: device1.Id,
		ScriptId: scripts["mb_dev1_turn_off_light3_v1"].Id,
	}
	ok, _ = deviceAction7.Valid()
	So(ok, ShouldEqual, true)
	deviceAction7.Id, err = adaptors.DeviceAction.Add(deviceAction7)
	So(err, ShouldBeNil)
	deviceActions["mb_dev1_turn_off_light3_v1"] = deviceAction7

	// light4
	deviceAction8 := &m.DeviceAction{
		Name:     "turn on light4",
		DeviceId: device1.Id,
		ScriptId: scripts["mb_dev1_turn_on_light4_v1"].Id,
	}
	ok, _ = deviceAction8.Valid()
	So(ok, ShouldEqual, true)
	deviceAction8.Id, err = adaptors.DeviceAction.Add(deviceAction8)
	So(err, ShouldBeNil)
	deviceActions["mb_dev1_turn_on_light4_v1"] = deviceAction8

	deviceAction9 := &m.DeviceAction{
		Name:     "turn off light4",
		DeviceId: device1.Id,
		ScriptId: scripts["mb_dev1_turn_off_light4_v1"].Id,
	}
	ok, _ = deviceAction9.Valid()
	So(ok, ShouldEqual, true)
	deviceAction9.Id, err = adaptors.DeviceAction.Add(deviceAction9)
	So(err, ShouldBeNil)
	deviceActions["mb_dev1_turn_off_light4_v1"] = deviceAction9

	// fan1
	deviceAction10 := &m.DeviceAction{
		Name:     "turn on fan1",
		DeviceId: device1.Id,
		ScriptId: scripts["mb_dev1_turn_on_fan1_v1"].Id,
	}
	ok, _ = deviceAction10.Valid()
	So(ok, ShouldEqual, true)
	deviceAction10.Id, err = adaptors.DeviceAction.Add(deviceAction10)
	So(err, ShouldBeNil)
	deviceActions["mb_dev1_turn_on_fan1_v1"] = deviceAction10

	deviceAction11 := &m.DeviceAction{
		Name:     "turn off fan1",
		DeviceId: device1.Id,
		ScriptId: scripts["mb_dev1_turn_off_fan1_v1"].Id,
	}
	ok, _ = deviceAction11.Valid()
	So(ok, ShouldEqual, true)
	deviceAction11.Id, err = adaptors.DeviceAction.Add(deviceAction11)
	So(err, ShouldBeNil)
	deviceActions["mb_dev1_turn_off_fan1_v1"] = deviceAction11

	// controll all lights
	deviceAction12 := &m.DeviceAction{
		Name:     "turn on all lights",
		DeviceId: device1.Id,
		ScriptId: scripts["mb_dev1_turn_on_all_lights_v1"].Id,
	}
	ok, _ = deviceAction12.Valid()
	So(ok, ShouldEqual, true)
	deviceAction12.Id, err = adaptors.DeviceAction.Add(deviceAction12)
	So(err, ShouldBeNil)
	deviceActions["mb_dev1_turn_on_all_lights_v1"] = deviceAction12

	deviceAction13 := &m.DeviceAction{
		Name:     "turn off all lights",
		DeviceId: device1.Id,
		ScriptId: scripts["mb_dev1_turn_off_all_lights_v1"].Id,
	}
	ok, _ = deviceAction13.Valid()
	So(ok, ShouldEqual, true)
	deviceAction13.Id, err = adaptors.DeviceAction.Add(deviceAction13)
	So(err, ShouldBeNil)
	deviceActions["mb_dev1_turn_off_all_lights_v1"] = deviceAction13

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
	stateDev1Temp1On := &m.DeviceState{
		SystemName:  "TEMP_1_ON",
		Description: "device temp 1 on",
		DeviceId:    device1.Id,
	}
	stateDev1Temp1Off := &m.DeviceState{
		SystemName:  "TEMP_1_OFF",
		Description: "device temp 1 off",
		DeviceId:    device1.Id,
	}
	stateDev1Temp2On := &m.DeviceState{
		SystemName:  "TEMP_2_ON",
		Description: "device temp 2 on",
		DeviceId:    device1.Id,
	}
	stateDev1Temp2Off := &m.DeviceState{
		SystemName:  "TEMP_2_OFF",
		Description: "device temp 2 off",
		DeviceId:    device1.Id,
	}
	stateMainGateOpened := &m.DeviceState{
		SystemName:  "DOOR_MAIN_OPENED",
		Description: "door opened",
		DeviceId:    device1.Id,
	}
	stateMainGateClosed := &m.DeviceState{
		SystemName:  "DOOR_MAIN_CLOSED",
		Description: "door closed",
		DeviceId:    device1.Id,
	}
	stateSecondGateOpened := &m.DeviceState{
		SystemName:  "DOOR_SECOND_OPENED",
		Description: "door opened",
		DeviceId:    device1.Id,
	}
	stateSecondGateClosed := &m.DeviceState{
		SystemName:  "DOOR_SECOND_CLOSED",
		Description: "door closed",
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
	ok, _ = stateDev1Temp1On.Valid()
	So(ok, ShouldEqual, true)
	ok, _ = stateDev1Temp1Off.Valid()
	So(ok, ShouldEqual, true)
	ok, _ = stateDev1Temp2On.Valid()
	So(ok, ShouldEqual, true)
	ok, _ = stateDev1Temp2Off.Valid()
	So(ok, ShouldEqual, true)
	ok, _ = stateMainGateOpened.Valid()
	So(ok, ShouldEqual, true)
	ok, _ = stateMainGateClosed.Valid()
	So(ok, ShouldEqual, true)
	ok, _ = stateSecondGateOpened.Valid()
	So(ok, ShouldEqual, true)
	ok, _ = stateSecondGateClosed.Valid()
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
	stateDev1Temp1On.Id, err = adaptors.DeviceState.Add(stateDev1Temp1On)
	So(err, ShouldBeNil)
	stateDev1Temp1Off.Id, err = adaptors.DeviceState.Add(stateDev1Temp1Off)
	So(err, ShouldBeNil)
	stateDev1Temp2On.Id, err = adaptors.DeviceState.Add(stateDev1Temp2On)
	So(err, ShouldBeNil)
	stateDev1Temp2Off.Id, err = adaptors.DeviceState.Add(stateDev1Temp2Off)
	So(err, ShouldBeNil)
	stateMainGateOpened.Id, err = adaptors.DeviceState.Add(stateMainGateOpened)
	So(err, ShouldBeNil)
	stateMainGateClosed.Id, err = adaptors.DeviceState.Add(stateMainGateClosed)
	So(err, ShouldBeNil)
	stateSecondGateOpened.Id, err = adaptors.DeviceState.Add(stateSecondGateOpened)
	So(err, ShouldBeNil)
	stateSecondGateClosed.Id, err = adaptors.DeviceState.Add(stateSecondGateClosed)
	So(err, ShouldBeNil)

	deviceStates["dev1_enabled"] = stateDev1Enabled
	deviceStates["dev1_disabled"] = stateDev1Disabled
	deviceStates["dev1_error"] = stateDev1Error
	deviceStates["dev1_light1_on"] = stateDev1Light1On
	deviceStates["dev1_light1_off"] = stateDev1Light1Off
	deviceStates["dev1_light2_on"] = stateDev1Light2On
	deviceStates["dev1_light2_off"] = stateDev1Light2Off
	deviceStates["dev1_light3_on"] = stateDev1Light3On
	deviceStates["dev1_light3_off"] = stateDev1Light3Off
	deviceStates["dev1_light4_on"] = stateDev1Light4On
	deviceStates["dev1_light4_off"] = stateDev1Light4Off
	deviceStates["dev1_fan1_on"] = stateDev1Fan1On
	deviceStates["dev1_fan1_off"] = stateDev1Fan1Off
	deviceStates["dev1_temp1_on"] = stateDev1Temp1On
	deviceStates["dev1_temp1_off"] = stateDev1Temp1Off
	deviceStates["dev1_temp2_on"] = stateDev1Temp2On
	deviceStates["dev1_temp2_off"] = stateDev1Temp2Off
	deviceStates["state_main_gate_opened"] = stateMainGateOpened
	deviceStates["state_main_gate_closed"] = stateMainGateClosed
	deviceStates["state_second_gate_opened"] = stateSecondGateOpened
	deviceStates["state_second_gate_closed"] = stateSecondGateClosed

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

	deviceActions["cmd_condition_check_v1"] = deviceAction21

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
		Type:       DevTypeModbusRtu,
		Node:       node1,
		IsGroup:    true,
		Properties: []byte("{}"),
	}

	modBusConfigDev3 := &DevModBusRtuConfig{
		Baud:     115200,
		DataBits: 8,
		StopBits: 2,
		Parity:   "none",
	}

	ok, _ = device3.SetProperties(modBusConfigDev3)
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
		Type:       DevTypeModbusRtu,
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
