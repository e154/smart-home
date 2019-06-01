package env1

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/system/initial/assertions"
	. "github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/null"
)

func addMaps(adaptors *adaptors.Adaptors,
	scripts map[string]*m.Script,
	devices []*m.Device,
	imageList map[string]*m.Image,
	deviceActions map[string]*m.DeviceAction,
	deviceStates map[string]*m.DeviceState) (maps []*m.Map) {

	var err error

	// map 1
	// ------------------------------------------------
	map1 := &m.Map{
		Name:        "office1",
		Description: "офис на ул. Красный проспект, д.22",
		Options: m.MapOptions{
			Zoom:              1,
			ElementStateText:  false,
			ElementOptionText: false,
		},
	}
	ok, _ := map1.Valid()
	So(ok, ShouldEqual, true)
	map1.Id, err = adaptors.Map.Add(map1)
	So(err, ShouldBeNil)

	// background layer
	// ------------------------------------------------
	mapLayer1 := &m.MapLayer{
		Name:        "background",
		Status:      "enabled",
		Description: "фон",
		MapId:       map1.Id,
		Weight:      0,
	}
	ok, _ = mapLayer1.Valid()
	So(ok, ShouldEqual, true)
	mapLayer1.Id, err = adaptors.MapLayer.Add(mapLayer1)
	So(err, ShouldBeNil)

	// base layer
	// ------------------------------------------------
	mapLayer2 := &m.MapLayer{
		Name:        "base",
		Status:      "enabled",
		Description: "базовый слой",
		MapId:       map1.Id,
		Weight:      1,
	}
	ok, _ = mapLayer2.Valid()
	So(ok, ShouldEqual, true)
	mapLayer2.Id, err = adaptors.MapLayer.Add(mapLayer2)
	So(err, ShouldBeNil)

	// light1
	// ------------------------------------------------
	devLight1 := &m.MapDevice{
		SystemName: "DEV1_LIGHT1",
		DeviceId:   devices[0].Id,
		ImageId:    imageList["lamp_v1_def"].Id,
		States: []*m.MapDeviceState{
			{
				DeviceStateId: deviceStates["dev1_light1_on"].Id,
				ImageId:       imageList["lamp_v1_y"].Id,
			},
			{
				DeviceStateId: deviceStates["dev1_light1_off"].Id,
				ImageId:       imageList["lamp_v1_def"].Id,
			},
			{
				DeviceStateId: deviceStates["dev1_error"].Id,
				ImageId:       imageList["lamp_v1_r"].Id,
			},
		},
		Actions: []*m.MapDeviceAction{
			{
				DeviceActionId: deviceActions["mb_dev1_turn_on_light1_v1"].Id,
				ImageId:        imageList["button_v1_on"].Id,
			},
			{
				DeviceActionId: deviceActions["mb_dev1_turn_off_light1_v1"].Id,
				ImageId:        imageList["button_v1_off"].Id,
			},
		},
	}

	ok, _ = devLight1.Valid()
	mapElementLight1 := &m.MapElement{
		Name: "dev1_light1",
		Prototype: m.Prototype{
			MapDevice: devLight1,
		},
		MapId:   map1.Id,
		LayerId: mapLayer2.Id,
		Status:  Enabled,
		GraphSettings: m.MapElementGraphSettings{
			Width:  null.NewInt64(33),
			Height: null.NewInt64(33),
			Position: m.MapElementGraphSettingsPosition{
				Top:  0,
				Left: 0,
			},
		},
	}
	ok, _ = mapElementLight1.Valid()
	So(ok, ShouldEqual, true)
	mapElementLight1.Id, err = adaptors.MapElement.Add(mapElementLight1)
	So(err, ShouldBeNil)

	// light2
	// ------------------------------------------------
	devLight2 := &m.MapDevice{
		SystemName: "DEV1_LIGHT2",
		DeviceId:   devices[0].Id,
		ImageId:    imageList["lamp_v1_def"].Id,
		States: []*m.MapDeviceState{
			{
				DeviceStateId: deviceStates["dev1_light2_on"].Id,
				ImageId:       imageList["lamp_v1_y"].Id,
			},
			{
				DeviceStateId: deviceStates["dev1_light2_off"].Id,
				ImageId:       imageList["lamp_v1_def"].Id,
			},
			{
				DeviceStateId: deviceStates["dev1_error"].Id,
				ImageId:       imageList["lamp_v1_r"].Id,
			},
		},
		Actions: []*m.MapDeviceAction{
			{
				DeviceActionId: deviceActions["mb_dev1_turn_on_light2_v1"].Id,
				ImageId:        imageList["button_v1_on"].Id,
			},
			{
				DeviceActionId: deviceActions["mb_dev1_turn_off_light2_v1"].Id,
				ImageId:        imageList["button_v1_off"].Id,
			},
		},
	}
	mapElementLight2 := &m.MapElement{
		Name: "dev1_light2",
		Prototype: m.Prototype{
			MapDevice: devLight2,
		},
		MapId:   map1.Id,
		LayerId: mapLayer2.Id,
		Status:  Enabled,
		GraphSettings: m.MapElementGraphSettings{
			Width:  null.NewInt64(33),
			Height: null.NewInt64(33),
			Position: m.MapElementGraphSettingsPosition{
				Top:  0,
				Left: 0,
			},
		},
	}
	ok, _ = mapElementLight2.Valid()
	So(ok, ShouldEqual, true)
	mapElementLight2.Id, err = adaptors.MapElement.Add(mapElementLight2)
	So(err, ShouldBeNil)

	// light3
	// ------------------------------------------------
	devLight3 := &m.MapDevice{
		SystemName: "DEV1_LIGHT3",
		DeviceId:   devices[0].Id,
		ImageId:    imageList["lamp_v1_def"].Id,
		States: []*m.MapDeviceState{
			{
				DeviceStateId: deviceStates["dev1_light3_on"].Id,
				ImageId:       imageList["lamp_v1_y"].Id,
			},
			{
				DeviceStateId: deviceStates["dev1_light3_off"].Id,
				ImageId:       imageList["lamp_v1_def"].Id,
			},
			{
				DeviceStateId: deviceStates["dev1_error"].Id,
				ImageId:       imageList["lamp_v1_r"].Id,
			},
		},
		Actions: []*m.MapDeviceAction{
			{
				DeviceActionId: deviceActions["mb_dev1_turn_on_light3_v1"].Id,
				ImageId:        imageList["button_v1_on"].Id,
			},
			{
				DeviceActionId: deviceActions["mb_dev1_turn_off_light3_v1"].Id,
				ImageId:        imageList["button_v1_off"].Id,
			},
		},
	}
	mapElementLight3 := &m.MapElement{
		Name: "dev1_light3",
		Prototype: m.Prototype{
			MapDevice: devLight3,
		},
		MapId:   map1.Id,
		LayerId: mapLayer2.Id,
		Status:  Enabled,
		GraphSettings: m.MapElementGraphSettings{
			Width:  null.NewInt64(33),
			Height: null.NewInt64(33),
			Position: m.MapElementGraphSettingsPosition{
				Top:  0,
				Left: 0,
			},
		},
	}
	ok, _ = mapElementLight3.Valid()
	So(ok, ShouldEqual, true)
	mapElementLight3.Id, err = adaptors.MapElement.Add(mapElementLight3)
	So(err, ShouldBeNil)

	// light4
	// ------------------------------------------------
	devLight4 := &m.MapDevice{
		SystemName: "DEV1_LIGHT4",
		DeviceId:   devices[0].Id,
		ImageId:    imageList["lamp_v1_def"].Id,
		States: []*m.MapDeviceState{
			{
				DeviceStateId: deviceStates["dev1_light4_on"].Id,
				ImageId:       imageList["lamp_v1_y"].Id,
			},
			{
				DeviceStateId: deviceStates["dev1_light4_off"].Id,
				ImageId:       imageList["lamp_v1_def"].Id,
			},
			{
				DeviceStateId: deviceStates["dev1_error"].Id,
				ImageId:       imageList["lamp_v1_r"].Id,
			},
		},
		Actions: []*m.MapDeviceAction{
			{
				DeviceActionId: deviceActions["mb_dev1_turn_on_light4_v1"].Id,
				ImageId:        imageList["button_v1_on"].Id,
			},
			{
				DeviceActionId: deviceActions["mb_dev1_turn_off_light4_v1"].Id,
				ImageId:        imageList["button_v1_off"].Id,
			},
		},
	}
	mapElementLight4 := &m.MapElement{
		Name: "dev1_light4",
		Prototype: m.Prototype{
			MapDevice: devLight4,
		},
		MapId:   map1.Id,
		LayerId: mapLayer2.Id,
		Status:  Enabled,
		GraphSettings: m.MapElementGraphSettings{
			Width:  null.NewInt64(33),
			Height: null.NewInt64(33),
			Position: m.MapElementGraphSettingsPosition{
				Top:  0,
				Left: 0,
			},
		},
	}
	ok, _ = mapElementLight4.Valid()
	So(ok, ShouldEqual, true)
	mapElementLight4.Id, err = adaptors.MapElement.Add(mapElementLight4)
	So(err, ShouldBeNil)

	// fan5
	// ------------------------------------------------
	devFan1 := &m.MapDevice{
		SystemName: "DEV1_FAN1",
		DeviceId:   devices[0].Id,
		ImageId:    imageList["fan_v1_def"].Id,
		States: []*m.MapDeviceState{
			{
				DeviceStateId: deviceStates["dev1_fan1_on"].Id,
				ImageId:       imageList["fan_v1_y"].Id,
			},
			{
				DeviceStateId: deviceStates["dev1_fan1_off"].Id,
				ImageId:       imageList["fan_v1_def"].Id,
			},
			{
				DeviceStateId: deviceStates["dev1_error"].Id,
				ImageId:       imageList["fan_v1_r"].Id,
			},
		},
		Actions: []*m.MapDeviceAction{
			{
				DeviceActionId: deviceActions["mb_dev1_turn_on_fan1_v1"].Id,
				ImageId:        imageList["button_v1_on"].Id,
			},
			{
				DeviceActionId: deviceActions["mb_dev1_turn_off_fan1_v1"].Id,
				ImageId:        imageList["button_v1_off"].Id,
			},
		},
	}
	mapElementFan1 := &m.MapElement{
		Name: "dev1_fan1",
		Prototype: m.Prototype{
			MapDevice: devFan1,
		},
		MapId:   map1.Id,
		LayerId: mapLayer2.Id,
		Status:  Enabled,
		GraphSettings: m.MapElementGraphSettings{
			Width:  null.NewInt64(33),
			Height: null.NewInt64(33),
			Position: m.MapElementGraphSettingsPosition{
				Top:  0,
				Left: 0,
			},
		},
	}
	ok, _ = mapElementFan1.Valid()
	So(ok, ShouldEqual, true)
	mapElementFan1.Id, err = adaptors.MapElement.Add(mapElementFan1)
	So(err, ShouldBeNil)

	// temp1
	// ------------------------------------------------
	dev1Temp1 := &m.MapDevice{
		SystemName: "DEV1_TEMP1",
		DeviceId:   devices[0].Id,
		ImageId:    imageList["temp_v1_def"].Id,
		States: []*m.MapDeviceState{
			{
				DeviceStateId: deviceStates["dev1_temp1_on"].Id,
				ImageId:       imageList["temp_v1_y"].Id,
			},
			{
				DeviceStateId: deviceStates["dev1_temp1_off"].Id,
				ImageId:       imageList["temp_v1_r"].Id,
			},
		},
	}
	mapElementTemp1 := &m.MapElement{
		Name:        "dev1_temp1",
		Description: "temperature sensor room1",
		Prototype: m.Prototype{
			MapDevice: dev1Temp1,
		},
		MapId:   map1.Id,
		LayerId: mapLayer2.Id,
		Status:  Enabled,
		GraphSettings: m.MapElementGraphSettings{
			Width:  null.NewInt64(33),
			Height: null.NewInt64(33),
			Position: m.MapElementGraphSettingsPosition{
				Top:  0,
				Left: 0,
			},
		},
	}
	ok, _ = mapElementTemp1.Valid()
	So(ok, ShouldEqual, true)
	mapElementTemp1.Id, err = adaptors.MapElement.Add(mapElementTemp1)
	So(err, ShouldBeNil)

	// temp2
	// ------------------------------------------------
	dev1Temp2 := &m.MapDevice{
		SystemName: "DEV1_TEMP2",
		DeviceId:   devices[0].Id,
		ImageId:    imageList["temp_v1_def"].Id,
		States: []*m.MapDeviceState{
			{
				DeviceStateId: deviceStates["dev1_temp2_on"].Id,
				ImageId:       imageList["temp_v1_y"].Id,
			},
			{
				DeviceStateId: deviceStates["dev1_temp2_off"].Id,
				ImageId:       imageList["temp_v1_r"].Id,
			},
		},
	}
	mapElementTemp2 := &m.MapElement{
		Name:        "dev1_temp2",
		Description: "temperature sensor room2",
		Prototype: m.Prototype{
			MapDevice: dev1Temp2,
		},
		MapId:   map1.Id,
		LayerId: mapLayer2.Id,
		Status:  Enabled,
		GraphSettings: m.MapElementGraphSettings{
			Width:  null.NewInt64(33),
			Height: null.NewInt64(33),
			Position: m.MapElementGraphSettingsPosition{
				Top:  0,
				Left: 0,
			},
		},
	}
	ok, _ = mapElementTemp2.Valid()
	So(ok, ShouldEqual, true)
	mapElementTemp2.Id, err = adaptors.MapElement.Add(mapElementTemp2)
	So(err, ShouldBeNil)

	// map element text1
	// ------------------------------------------------
	mapText1 := &m.MapText{
		Text: "background",
	}

	mapElementText1 := &m.MapElement{
		Name: "text1",
		Prototype: m.Prototype{
			MapText: mapText1,
		},
		MapId:   map1.Id,
		LayerId: mapLayer2.Id,
		Status:  Enabled,
		GraphSettings: m.MapElementGraphSettings{
			Width:  null.NewInt64(33),
			Height: null.NewInt64(33),
			Position: m.MapElementGraphSettingsPosition{
				Top:  0,
				Left: 0,
			},
		},
	}

	ok, _ = mapElementText1.Valid()
	So(ok, ShouldEqual, true)
	mapElementText1.Id, err = adaptors.MapElement.Add(mapElementText1)
	So(err, ShouldBeNil)

	return
}
