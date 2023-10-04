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

package events

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEventName(t *testing.T) {

	event1 := EventStateChanged{}
	event2 := &EventStateChanged{}
	event3 := &EventEntityUnloaded{}

	require.Equal(t, EventName(event1), "event_state_changed")
	require.Equal(t, EventName(event2), "event_state_changed")
	require.Equal(t, EventName(event3), "event_entity_unloaded")
}
