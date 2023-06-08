package stream

import (
	"encoding/json"

	"github.com/e154/smart-home/plugins/webpush"
	"github.com/e154/smart-home/common/events"
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
	case webpush.EventNewWebPushPublicKey:
		go e.eventNewWebPushPublicKey(v)
	case events.EventDirectMessage:
		go e.eventDirectMessage(v.UserID, v.Query, v.Message)
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
	body, _ := json.Marshal(msg)
	e.directMessage(userID, query, body)
}
