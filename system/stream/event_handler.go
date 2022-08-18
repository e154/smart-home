package stream

import (
	"encoding/json"

	events2 "github.com/e154/smart-home/common/events"
)

type eventHandler struct {
	broadcast func(query string, message []byte)
}

func NewEventHandler(broadcast func(query string, message []byte)) *eventHandler {
	return &eventHandler{broadcast: broadcast}
}

func (e *eventHandler) eventHandler(_ string, message interface{}) {

	switch msg := message.(type) {
	case events2.EventStateChanged:
		go e.eventStateChangedHandler(msg)
	case events2.EventLoadedPlugin:
	case events2.EventUnloadedPlugin:
	case events2.EventCreatedEntity:
	case events2.EventUpdatedEntity:
	case events2.EventDeletedEntity:
	case events2.EventEntitySetState:
	}
}

func (e *eventHandler) eventStateChangedHandler(msg events2.EventStateChanged) {
	//todo optimize
	b, _ := json.Marshal(msg)
	e.broadcast("state_changed", b)
}
