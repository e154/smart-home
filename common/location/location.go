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

package location

import (
	"encoding/json"
	"fmt"

	"github.com/e154/smart-home/common/web"
	m "github.com/e154/smart-home/models"
)

const (
	// IpApi ...
	IpApi = "http://ip-api.com/json"
	// IPAPI ...
	IPAPI = "https://ipapi.co/json/"
)

// GeoLocationFromIP ...
func GeoLocationFromIP(ip string) (location m.GeoLocation, err error) {

	var body []byte
	if _, body, err = web.Probe(web.Request{Method: "GET", Url: fmt.Sprintf("%s/%s", IpApi, ip)}); err != nil {
		return
	}
	location = m.GeoLocation{}
	err = json.Unmarshal(body, &location)

	return
}

// GetRegionInfo ...
func GetRegionInfo() (info m.RegionInfo, err error) {

	var body []byte
	if _, body, err = web.Probe(web.Request{Method: "GET", Url: IPAPI}); err != nil {
		return
	}
	info = m.RegionInfo{}
	err = json.Unmarshal(body, &info)

	return
}
