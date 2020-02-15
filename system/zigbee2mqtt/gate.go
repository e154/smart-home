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

package zigbee2mqtt

const (
	baseTopic = "zigbee2mqtt/bridge"
)

type Gate struct {

}

func (g *Gate) getState() {}
func (g *Gate) getConfig() {}
func (g *Gate) getLog() {}
func (g *Gate) getDevices() {}
func (g *Gate) configPermitJoin() {}
func (g *Gate) configLastSeen() {}
func (g *Gate) configElapsed() {}
func (g *Gate) configReset() {}
func (g *Gate) configLogLevel() {}
func (g *Gate) DeviceOptions() {}
func (g *Gate) Remove() {}
func (g *Gate) Ban() {}
func (g *Gate) Whitelist() {}
func (g *Gate) Rename() {}
func (g *Gate) RenameLast() {}
func (g *Gate) AddGroup() {}
func (g *Gate) RemoveGroup() {}
func (g *Gate) Networkmap() {}