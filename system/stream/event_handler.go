package stream

import (
	"encoding/json"

	"github.com/e154/smart-home/common/events"
)

type eventHandler struct {
	broadcast func(query string, message []byte)
}

func NewEventHandler(broadcast func(query string, message []byte)) *eventHandler {
	return &eventHandler{broadcast: broadcast}
}

func (e *eventHandler) eventHandler(_ string, message interface{}) {

	switch message.(type) {
	case events.EventStateChanged:
		go e.eventStateChangedHandler(message)
	case events.EventLastStateChanged:
		go e.eventStateChangedHandler(message)
	case events.EventLoadedPlugin:
	case events.EventUnloadedPlugin:
	case events.EventCreatedEntity:
	case events.EventUpdatedEntity:
	case events.EventDeletedEntity:
	case events.EventEntitySetState:
	}
}

func (e *eventHandler) eventStateChangedHandler(msg interface{}) {
	//todo optimize
	b, _ := json.Marshal(msg)
	e.broadcast("state_changed", b)
}
