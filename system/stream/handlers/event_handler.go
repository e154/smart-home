package handlers

import (
	"context"
	"encoding/json"

	"go.uber.org/fx"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/webpush"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/stream"
)

type EventHandler struct {
	stream   *stream.Stream
	eventBus bus.Bus
}

func NewEventHandler(lc fx.Lifecycle,
	stream *stream.Stream,
	eventBus bus.Bus) *EventHandler {
	handler := &EventHandler{
		stream:   stream,
		eventBus: eventBus,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			return handler.Start(ctx)
		},
		OnStop: func(ctx context.Context) (err error) {
			return handler.Shutdown(ctx)
		},
	})

	return handler
}

// Start ...
func (s *EventHandler) Start(_ context.Context) error {
	s.stream.Subscribe("event_get_last_state", s.EventGetLastState)
	s.stream.Subscribe("event_update_device_location", s.EventUpdateDeviceLocation)
	s.stream.Subscribe("event_add_webpush_subscription", s.EventAddWebPushSubscription)
	s.stream.Subscribe("event_get_webpush_public_key", s.EventGetWebPushPublicKey)
	return nil
}

// Shutdown ...
func (s *EventHandler) Shutdown(_ context.Context) error {
	s.stream.UnSubscribe("event_get_last_state")
	s.stream.UnSubscribe("event_update_device_location")
	s.stream.UnSubscribe("event_add_webpush_subscription")
	s.stream.UnSubscribe("event_get_webpush_public_key")
	return nil
}

func (s *EventHandler) EventUpdateDeviceLocation(client stream.IStreamClient, query string, body []byte) {
	//var userID int64
	//if user := client.GetUser(); user != nil {
	//	userID = user.Id
	//}
	//fmt.Println(userID, string(body))
}

func (s *EventHandler) EventGetWebPushPublicKey(client stream.IStreamClient, query string, body []byte) {
	var userID int64
	if user := client.GetUser(); user != nil {
		userID = user.Id
	}

	s.eventBus.Publish(webpush.TopicPluginWebpush, webpush.EventGetWebPushPublicKey{
		UserID: userID,
	})
}

func (s *EventHandler) EventAddWebPushSubscription(client stream.IStreamClient, query string, body []byte) {
	var userID int64
	if user := client.GetUser(); user != nil {
		userID = user.Id
	}

	subscription := &m.Subscription{}
	_ = json.Unmarshal(body, subscription)
	s.eventBus.Publish(webpush.TopicPluginWebpush, webpush.EventAddWebPushSubscription{
		UserID:       userID,
		Subscription: subscription,
	})
}

func (s *EventHandler) EventGetLastState(client stream.IStreamClient, query string, body []byte) {
	req := map[string]common.EntityId{}
	_ = json.Unmarshal(body, &req)
	id := req["entity_id"]
	s.eventBus.Publish("system/entities/"+id.String(), events.EventGetLastState{
		EntityId: id,
	})
}
