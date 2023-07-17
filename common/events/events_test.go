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
	require.Equal(t, EventName(event3), "event_entity_deleted")
}
