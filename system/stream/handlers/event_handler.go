package handlers

import (
	"context"
	"encoding/json"

	"go.uber.org/fx"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/events"
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
	return nil
}

// Shutdown ...
func (s *EventHandler) Shutdown(_ context.Context) error {
	s.stream.UnSubscribe("event_get_last_state")
	return nil
}

func (s *EventHandler) EventGetLastState(client stream.IStreamClient, query string, body []byte) {
	req := map[string]common.EntityId{}
	_ = json.Unmarshal(body, &req)
	s.eventBus.Publish(bus.TopicEntities, events.EventGetLastState{
		EntityId: req["entity_id"],
	})
}
