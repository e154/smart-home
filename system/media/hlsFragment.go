// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

package media

import (
	"time"

	"github.com/deepch/vdk/av"
)

// Fragment struct
type Fragment struct {
	Independent bool          //Fragment have i-frame (key frame)
	Finish      bool          //Fragment Ready
	Duration    time.Duration //Fragment Duration
	Packets     []*av.Packet  //Packet Slice
}

// NewFragment open new fragment
func (element *Segment) NewFragment() *Fragment {
	res := &Fragment{}
	element.Fragment[element.CurrentFragmentID] = res
	return res
}

// GetDuration return fragment dur
func (element *Fragment) GetDuration() time.Duration {
	return element.Duration
}

// WritePacket to fragment func
func (element *Fragment) WritePacket(packet *av.Packet) {
	//increase fragment dur
	element.Duration += packet.Duration
	//Independent if have key
	if packet.IsKeyFrame {
		element.Independent = true
	}
	//append packet to slice of packet
	element.Packets = append(element.Packets, packet)
}

// Close fragment block func
func (element *Fragment) Close() {
	//TODO add callback func
	//finalize fragment
	element.Finish = true
}
