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

package bus

import (
	"testing"
)

func TestTopicMatch(t *testing.T) {

	const topic = "myhome/groundfloor/livingroom/temperature"
	const topic2 = "myhome/groundfloor/kitchen/temperature"

	// Test cases with exact matches
	assertMatch(t, []byte(topic), []byte("myhome/groundfloor/livingroom/temperature"))

	// Test cases with wildcard matches

	assertMatch(t, []byte(topic), []byte("myhome/groundfloor/#"))
	assertMatch(t, []byte(topic2), []byte("myhome/groundfloor/#"))
	assertMatch(t, []byte(topic), []byte("myhome/groundfloor/+/#"))
	assertMatch(t, []byte(topic2), []byte("myhome/groundfloor/+/#"))

	// Test cases with no matches
	assertNoMatch(t, []byte(topic), []byte("myhome/groundfloor/livingroom/temperature/"))
	assertNoMatch(t, []byte(topic), []byte("myhome/groundfloor/livingroom/"))
	assertNoMatch(t, []byte(topic), []byte("myhome/groundfloor/+/temperature/"))
	assertNoMatch(t, []byte("myhome/groundfloor/livingroom"), []byte("myhome/groundfloor/+/temperature"))
}

func assertMatch(t *testing.T, topic []byte, topicFilter []byte) {
	if !TopicMatch(topic, topicFilter) {
		t.Errorf("Expected topic %s to match filter %s", topic, topicFilter)
	}
}

func assertNoMatch(t *testing.T, topic []byte, topicFilter []byte) {
	if TopicMatch(topic, topicFilter) {
		t.Errorf("Expected topic %s to not match filter %s", topic, topicFilter)
	}
}
