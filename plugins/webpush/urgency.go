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

package webpush

// Urgency indicates to the push service how important a message is to the user.
// This can be used by the push service to help conserve the battery life of a user's device
// by only waking up for important messages when battery is low.
type Urgency string

const (
	// UrgencyVeryLow requires device state: on power and Wi-Fi
	UrgencyVeryLow Urgency = "very-low"
	// UrgencyLow requires device state: on either power or Wi-Fi
	UrgencyLow Urgency = "low"
	// UrgencyNormal excludes device state: low battery
	UrgencyNormal Urgency = "normal"
	// UrgencyHigh admits device state: low battery
	UrgencyHigh Urgency = "high"
)

// Checking allowable values for the urgency header
func isValidUrgency(urgency Urgency) bool {
	switch urgency {
	case UrgencyVeryLow, UrgencyLow, UrgencyNormal, UrgencyHigh:
		return true
	}
	return false
}
