package media

import (
	"context"
	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/common/logger"
	"github.com/e154/smart-home/plugins/onvif"
	"github.com/e154/smart-home/system/bus"
	"go.uber.org/fx"
)

var (
	log = logger.MustGetLogger("media")
)

type Media struct {
	storage  *StorageST
	eventBus bus.Bus
}

func NewMedia(lc fx.Lifecycle,
	eventBus bus.Bus) *Media {
	rtsp := &Media{
		eventBus: eventBus,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			return rtsp.Start(ctx)
		},
		OnStop: func(ctx context.Context) (err error) {
			return rtsp.Shutdown(ctx)
		},
	})

	return rtsp
}

func (r *Media) Start(ctx context.Context) (err error) {
	go RTSPServer()
	go Storage.StreamChannelRunAll()

	_ = r.eventBus.Subscribe("system/entities/+", r.eventHandler)

	return
}

func (r *Media) Shutdown(ctx context.Context) (err error) {
	_ = r.eventBus.Unsubscribe("system/entities/+", r.eventHandler)

	Storage.StopAll()
	return
}

// eventHandler ...
func (r *Media) eventHandler(_ string, message interface{}) {

	switch event := message.(type) {
	case events.EventEntityUnloaded:
		go r.eventEntityUnloaded(event)
	case events.EventEntityLoaded:
		go r.eventEntityLoaded(event)
	case events.EventStateChanged:
		go r.eventStateChanged(event)
	}
}

func (r *Media) eventEntityUnloaded(event events.EventEntityUnloaded) {
	if event.PluginName != onvif.Name {
		return
	}
	if err := Storage.StreamDelete(event.EntityId.String()); err != nil {
		log.Error(err.Error())
	}
}

func (r *Media) eventEntityLoaded(event events.EventEntityLoaded) {
	if event.PluginName != onvif.Name {
		return
	}
}

func (r *Media) eventStateChanged(event events.EventStateChanged) {
	if event.PluginName != onvif.Name {
		return
	}
	streamUri := event.NewState.Attributes[onvif.AttrStreamUri].Decrypt()
	if streamUri == "" {
		return
	}

	// add/update stream
	payload := StreamST{
		Name: event.EntityId.String(),
		Channels: map[string]ChannelST{
			"0": {
				URL:      streamUri,
				OnDemand: true,
			},
		},
	}
	if err := Storage.StreamAdd(event.EntityId.String(), payload); err != nil {
		if err = Storage.StreamEdit(event.EntityId.String(), payload); err != nil {
			log.Error(err.Error())
		}
	}
}
