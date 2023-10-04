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

// Segment struct
type Segment struct {
	FPS               int               //Current fps
	CurrentFragment   *Fragment         //CurrentFragment link
	CurrentFragmentID int               //CurrentFragment ID
	Finish            bool              //Segment Ready
	Duration          time.Duration     //Segment Duration
	Time              time.Time         //Realtime EXT-X-PROGRAM-DATE-TIME
	Fragment          map[int]*Fragment //Fragment map
}

// NewSegment func
func (element *MuxerHLS) NewSegment() *Segment {
	res := &Segment{
		Fragment:          make(map[int]*Fragment),
		CurrentFragmentID: -1, //Default fragment -1
	}
	//Increase MSN
	element.MSN++
	element.Segments[element.MSN] = res
	return res
}

// GetDuration func
func (element *Segment) GetDuration() time.Duration {
	return element.Duration
}

// SetFPS func
func (element *Segment) SetFPS(fps int) {
	element.FPS = fps
}

// WritePacket func
func (element *Segment) WritePacket(packet *av.Packet) {
	if element.CurrentFragment == nil || element.CurrentFragment.GetDuration().Milliseconds() >= element.FragmentMS(element.FPS) {
		if element.CurrentFragment != nil {
			element.CurrentFragment.Close()
		}
		element.CurrentFragmentID++
		element.CurrentFragment = element.NewFragment()
	}
	element.Duration += packet.Duration
	element.CurrentFragment.WritePacket(packet)
}

// GetFragmentID func
func (element *Segment) GetFragmentID() int {
	return element.CurrentFragmentID
}

// Close segment func
func (element *Segment) Close() {
	element.Finish = true
	if element.CurrentFragment != nil {
		element.CurrentFragment.Close()
	}
}

// FragmentMS func
func (element *Segment) FragmentMS(fps int) int64 {
	for i := 6; i >= 1; i-- {
		if fps%i == 0 {
			return int64(float64(1000) / float64(fps) * float64(i))
		}
	}
	return 100
}
