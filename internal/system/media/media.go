// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package media

import (
	"context"
	"fmt"

	"github.com/e154/smart-home/pkg/logger"

	"github.com/e154/bus"
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
	//go RTSPServer()
	go Storage.StreamChannelRunAll()
	_ = r.eventBus.Subscribe("system/media/#", r.eventHandler)

	return
}

func (r *Media) Shutdown(ctx context.Context) (err error) {
	_ = r.eventBus.Unsubscribe("system/media/#", r.eventHandler)

	Storage.StopAll()
	return
}

// eventHandler ...
func (r *Media) eventHandler(_ string, message interface{}) {

	switch event := message.(type) {
	case EventRemoveList:
		go r.eventRemoveList(event)
	case EventUpdateList:
		go r.eventUpdateList(event)
	}
}

func (r *Media) eventRemoveList(event EventRemoveList) {
	if event.Name == "" {
		return
	}
	if err := Storage.StreamDelete(event.Name); err != nil {
		log.Error(err.Error())
	}
}

func (r *Media) eventUpdateList(event EventUpdateList) {
	if event.Name == "" || len(event.Channels) == 0 {
		return
	}

	// add/update stream
	payload := StreamST{
		Name:     event.Name,
		Channels: make(map[string]ChannelST),
	}

	for i, item := range event.Channels {
		payload.Channels[fmt.Sprintf("%d", i)] = ChannelST{
			URL:      item,
			OnDemand: true,
		}
	}

	if err := Storage.StreamAdd(event.Name, payload); err != nil {
		if err = Storage.StreamEdit(event.Name, payload); err != nil {
			log.Error(err.Error())
		}
	}
}
