package stream

import (
	"encoding/json"
	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/plugins/webpush"
)

type eventHandler struct {
	broadcast     func(query string, message []byte)
	directMessage func(userID int64, query string, message []byte)
}

func NewEventHandler(broadcast func(query string, message []byte),
	directMessage func(userID int64, query string, message []byte)) *eventHandler {
	return &eventHandler{
		broadcast:     broadcast,
		directMessage: directMessage,
	}
}

func (e *eventHandler) eventHandler(_ string, message interface{}) {

	switch v := message.(type) {

	// entities
	case events.EventStateChanged:
		go e.eventStateChangedHandler(message)
	case events.EventLastStateChanged:
		go e.eventStateChangedHandler(message)
	case events.EventCreatedEntity:
	case events.EventUpdatedEntity:
	case events.CommandUnloadEntity:
	case events.EventEntityLoaded:
		go e.event(message)
	case events.EventEntityUnloaded:
		go e.event(message)
	case events.EventEntitySetState:

	// notifications
	case webpush.EventNewWebPushPublicKey:
		go e.eventNewWebPushPublicKey(v)
	case events.EventDirectMessage:
		go e.eventDirectMessage(v.UserID, v.Query, v.Message)

	// plugins
	case events.EventLoadedPlugin:
	case events.EventUnloadedPlugin:

	// tasks
	case events.EventTaskLoaded:
		go e.event(message)
	case events.EventTaskUnloaded:
		go e.event(message)
	case events.EventTaskCompleted:
		go e.event(message)

	// triggers
	case events.EventTriggerLoaded:
		go e.event(message)
	case events.EventTriggerUnloaded:
		go e.event(message)
	case events.EventTriggerCompleted:
		go e.event(message)

	// actions
	case events.EventActionCompleted:
		go e.event(message)
	}
}

func (e *eventHandler) eventNewWebPushPublicKey(event webpush.EventNewWebPushPublicKey) {
	b, _ := json.Marshal(event)
	if event.UserID != 0 {
		e.directMessage(event.UserID, "event_new_webpush_public_key", b)
		return
	}
	e.broadcast("event_new_webpush_public_key", b)
}

func (e *eventHandler) eventStateChangedHandler(msg interface{}) {
	//todo optimize
	b, _ := json.Marshal(msg)
	e.broadcast("state_changed", b)
}

func (e *eventHandler) eventDirectMessage(userID int64, query string, msg interface{}) {
	b, _ := json.Marshal(msg)
	e.directMessage(userID, query, b)
}

func (e *eventHandler) event(msg interface{}) {
	b, _ := json.Marshal(msg)
	e.broadcast(events.EventName(msg), b)
}
