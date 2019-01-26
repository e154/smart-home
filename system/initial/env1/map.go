package env1

import (
	"encoding/json"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/system/initial/assertions"
	. "github.com/e154/smart-home/common"
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
		Options:     json.RawMessage(`{"zoom":1,"element_state_text":false,"element_option_text":false}`),
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

	// device
	// ------------------------------------------------
	//light1
	devLight1 := &m.MapDevice{
		SystemName: "DEV1_LIGHT1",
		DeviceId:   devices[0].Id,
		States: []*m.MapDeviceState{
			{
				DeviceStateId: devices[0].Id,
				//ImageId: imageList[""]
			},
		},
		Actions: []*m.MapDeviceAction{},
	}

	ok, _ = devLight1.Valid()
	mapElementLight1 := &m.MapElement{
		Name:          "dev1_light1",
		Prototype:     devLight1,
		MapId:         map1.Id,
		LayerId:       mapLayer2.Id,
		Status:        Enabled,
		GraphSettings: json.RawMessage(`{"width":0,"height":0,"position":{"top":0,"left":0}}`),
	}
	ok, _ = mapElementLight1.Valid()
	So(ok, ShouldEqual, true)
	mapElementLight1.Id, err = adaptors.MapElement.Add(mapElementLight1)
	So(err, ShouldBeNil)

	//light2
	devLight2 := &m.MapDevice{
		SystemName: "DEV1_LIGHT2",
		DeviceId:   devices[0].Id,
	}
	mapElementLight2 := &m.MapElement{
		Name:          "dev1_light2",
		Prototype:     devLight2,
		MapId:         map1.Id,
		LayerId:       mapLayer2.Id,
		Status:        Enabled,
		GraphSettings: json.RawMessage(`{"width":0,"height":0,"position":{"top":0,"left":0}}`),
	}
	ok, _ = mapElementLight2.Valid()
	So(ok, ShouldEqual, true)
	mapElementLight2.Id, err = adaptors.MapElement.Add(mapElementLight2)
	So(err, ShouldBeNil)
	//light3
	devLight3 := &m.MapDevice{
		SystemName: "DEV1_LIGHT3",
		DeviceId:   devices[0].Id,
	}
	mapElementLight3 := &m.MapElement{
		Name:          "dev1_light3",
		Prototype:     devLight3,
		MapId:         map1.Id,
		LayerId:       mapLayer2.Id,
		Status:        Enabled,
		GraphSettings: json.RawMessage(`{"width":0,"height":0,"position":{"top":0,"left":0}}`),
	}
	ok, _ = mapElementLight3.Valid()
	So(ok, ShouldEqual, true)
	mapElementLight3.Id, err = adaptors.MapElement.Add(mapElementLight3)
	So(err, ShouldBeNil)
	//light4
	devLight4 := &m.MapDevice{
		SystemName: "DEV1_LIGHT4",
		DeviceId:   devices[0].Id,
	}
	mapElementLight4 := &m.MapElement{
		Name:          "dev1_light4",
		Prototype:     devLight4,
		MapId:         map1.Id,
		LayerId:       mapLayer2.Id,
		Status:        Enabled,
		GraphSettings: json.RawMessage(`{"width":0,"height":0,"position":{"top":0,"left":0}}`),
	}
	ok, _ = mapElementLight4.Valid()
	So(ok, ShouldEqual, true)
	mapElementLight4.Id, err = adaptors.MapElement.Add(mapElementLight4)
	So(err, ShouldBeNil)
	//fan5
	devFan1 := &m.MapDevice{
		SystemName: "DEV1_FAN1",
		DeviceId:   devices[0].Id,
	}
	mapElementFan1 := &m.MapElement{
		Name:          "dev1_fan5",
		Prototype:     devFan1,
		MapId:         map1.Id,
		LayerId:       mapLayer2.Id,
		Status:        Enabled,
		GraphSettings: json.RawMessage(`{"width":0,"height":0,"position":{"top":0,"left":0}}`),
	}
	ok, _ = mapElementFan1.Valid()
	So(ok, ShouldEqual, true)
	mapElementFan1.Id, err = adaptors.MapElement.Add(mapElementFan1)
	So(err, ShouldBeNil)

	// map element text1
	// ------------------------------------------------
	mapText1 := &m.MapText{
		Text: "background",
	}

	mapElementText1 := &m.MapElement{
		Name:          "text1",
		Prototype:     mapText1,
		MapId:         map1.Id,
		LayerId:       mapLayer2.Id,
		Status:        Enabled,
		GraphSettings: json.RawMessage(`{"width":0,"height":0,"position":{"top":0,"left":0}}`),
	}

	ok, _ = mapElementText1.Valid()
	So(ok, ShouldEqual, true)
	mapElementText1.Id, err = adaptors.MapElement.Add(mapElementText1)
	So(err, ShouldBeNil)

	return
}
