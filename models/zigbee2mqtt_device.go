// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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
	"fmt"
	"path/filepath"
	"time"
)

// Zigbee2mqttDevice ...
type Zigbee2mqttDevice struct {
	Id            string    `json:"id"`
	Zigbee2mqttId int64     `json:"zigbee2mqtt_id" validate:"required"`
	Name          string    `json:"name" validate:"required,max=254"`
	Type          string    `json:"type"`
	Model         string    `json:"model"`
	Description   string    `json:"description"`
	Manufacturer  string    `json:"manufacturer"`
	Functions     []string  `json:"functions"`
	ImageUrl      string    `json:"image_url"`
	Status        string    `json:"status"`
	Payload       []byte    `json:"payload"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// GetImageUrl ...
func (d *Zigbee2mqttDevice) GetImageUrl() {
	if d.Model != "" {
		d.ImageUrl = filepath.Join("/api_static", "devices", fmt.Sprintf("%s.jpg", d.Model))
	}
}
