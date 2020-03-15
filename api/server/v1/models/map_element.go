// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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

package models

import (
	"encoding/json"
	"time"
)

type MapElementGraphSettingsPosition struct {
	Top  int64 `json:"top"`
	Left int64 `json:"left"`
}
type MapElementGraphSettings struct {
	Width    *int64                          `json:"width"`
	Height   *int64                          `json:"height"`
	Position MapElementGraphSettingsPosition `json:"position"`
}

// swagger:model
type Prototype struct {
	*MapImage
	*MapText
	*MapDevice
}

func (n Prototype) MarshalJSON() (b []byte, err error) {
	switch {
	case n.MapText != nil:
		b, err = json.Marshal(n.MapText)
	case n.MapImage != nil:
		b, err = json.Marshal(n.MapImage)
	case n.MapDevice != nil && n.MapDevice.Device != nil:
		b, err = json.Marshal(n.MapDevice)
	default:
		b, err = json.Marshal(struct{}{})
		return
	}
	return
}

func (n *Prototype) UnmarshalJSON(data []byte) (err error) {

	device := &MapDevice{}
	err = json.Unmarshal(data, device)
	if device.Device != nil && device.Device.Id != 0 && device.DeviceId != 0 {
		n.MapDevice = device
		return
	}

	image := &MapImage{}
	err = json.Unmarshal(data, image)
	if image.ImageId != 0 {
		n.MapImage = image
		return
	}

	text := &MapText{}
	err = json.Unmarshal(data, text)
	n.MapText = text
	return
}

// swagger:model
type MapElement struct {
	Id            int64                   `json:"id"`
	Name          string                  `json:"name"`
	Description   string                  `json:"description"`
	PrototypeId   int64                   `json:"prototype_id"`
	PrototypeType string                  `json:"prototype_type"`
	Prototype     Prototype               `json:"prototype"`
	MapId         int64                   `json:"map_id"`
	LayerId       int64                   `json:"layer_id"`
	GraphSettings MapElementGraphSettings `json:"graph_settings"`
	Status        string                  `json:"status"`
	Weight        int                     `json:"weight"`
	Zone          *MapZone                `json:"zone"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`
}

type MapElementShort struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// swagger:model
type NewMapElement struct {
	Name          string                  `json:"name"`
	Description   string                  `json:"description"`
	PrototypeId   int64                   `json:"prototype_id"`
	PrototypeType string                  `json:"prototype_type"`
	Prototype     Prototype               `json:"prototype"`
	Map           *Map                    `json:"map"`
	Layer         *MapLayer               `json:"layer"`
	MapId         int64                   `json:"map_id"`
	LayerId       int64                   `json:"layer_id"`
	GraphSettings MapElementGraphSettings `json:"graph_settings"`
	Status        string                  `json:"status"`
	Weight        int                     `json:"weight"`
	Zone          *MapZone                `json:"zone"`
}

// swagger:model
type UpdateMapElement struct {
	Id            int64                   `json:"id"`
	Name          string                  `json:"name"`
	Description   string                  `json:"description"`
	PrototypeId   int64                   `json:"prototype_id"`
	PrototypeType string                  `json:"prototype_type"`
	Prototype     Prototype               `json:"prototype"`
	MapId         int64                   `json:"map_id"`
	LayerId       int64                   `json:"layer_id"`
	GraphSettings MapElementGraphSettings `json:"graph_settings"`
	Status        string                  `json:"status"`
	Weight        int                     `json:"weight"`
	Zone          *MapZone                `json:"zone"`
}

// swagger:model
type SortMapElement struct {
	Id     int64 `json:"id"`
	Weight int64 `json:"weight"`
}
