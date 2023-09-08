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
