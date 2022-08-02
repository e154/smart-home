package stream

import (
	"encoding/json"

	"github.com/e154/smart-home/system/event_bus/events"
)

type eventHandler struct {
	broadcast func(query string, message []byte)
}

func NewEventHandler(broadcast func(query string, message []byte)) *eventHandler {
	return &eventHandler{broadcast: broadcast}
}

func (e *eventHandler) eventHandler(_ string, message interface{}) {

	switch msg := message.(type) {
	case events.EventStateChanged:
		go e.eventStateChangedHandler(msg)
	case events.EventLoadedPlugin:
	case events.EventUnloadedPlugin:
	case events.EventCreatedEntity:
	case events.EventUpdatedEntity:
	case events.EventDeletedEntity:
	case events.EventEntitySetState:
	}
}

func (e *eventHandler) eventStateChangedHandler(msg events.EventStateChanged) {
	//todo optimize
	b, _ := json.Marshal(msg)
	e.broadcast("state_changed", b)
}
