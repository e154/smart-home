// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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

package pachka

import (
	"context"
	"fmt"
	"strconv"

	notify2 "github.com/e154/smart-home/internal/plugins/notify"
	notifyCommon "github.com/e154/smart-home/internal/plugins/notify/common"
	"github.com/e154/smart-home/internal/system/supervisor"
	"github.com/e154/smart-home/pkg/events"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
	"go.uber.org/atomic"
)

// Actor ...
type Actor struct {
	*supervisor.BaseActor
	isStarted   *atomic.Bool
	AccessToken string
	actionPool  chan events.EventCallEntityAction
	notify      *notify2.Notify
	client      *Client
}

// NewActor ...
func NewActor(entity *m.Entity,
	service plugins.Service) (*Actor, error) {

	settings := NewSettings()
	_, _ = settings.Deserialize(entity.Settings.Serialize())

	actor := &Actor{
		BaseActor:   supervisor.NewBaseActor(entity, service),
		actionPool:  make(chan events.EventCallEntityAction, 1000),
		isStarted:   atomic.NewBool(false),
		AccessToken: settings[AttrToken].Decrypt(),
		notify:      notify2.NewNotify(service.Adaptors()),
	}

	if actor.Attrs == nil {
		actor.Attrs = NewAttr()
	}

	if actor.Setts == nil {
		actor.Setts = NewSettings()
	}

	return actor, nil
}

func (e *Actor) Destroy() {
	if !e.isStarted.Load() {
		return
	}
	_ = e.Service.EventBus().Unsubscribe(notify2.TopicNotify, e.eventHandler)
	e.notify.Shutdown()

	if e.client != nil {
		e.client = nil
	}
	e.isStarted.Store(false)
}

func (e *Actor) Spawn() {

	var err error
	if e.isStarted.Load() {
		return
	}
	defer func() {
		if err == nil {
			e.isStarted.Store(true)
		}
	}()

	// load settings
	if e.Setts[AttrToken].Decrypt() != "" {
		e.client = NewClient(e.Setts[AttrToken].Decrypt())
	} else {
		log.Warn("empty access token")
		e.client = NewClient("NoToken")
	}

	_ = e.Service.EventBus().Subscribe(notify2.TopicNotify, e.eventHandler, false)
	e.notify.Start()

	e.BaseActor.Spawn()
}

// UpdateStatus ...
func (e *Actor) UpdateStatus() (err error) {

	var attributeValues = make(m.AttributeValue)
	// ...

	e.AttrMu.Lock()
	var changed bool
	if changed, err = e.Attrs.Deserialize(attributeValues); !changed {
		if err != nil {
			log.Warn(err.Error())
		}
	}
	e.AttrMu.Unlock()

	e.SaveState(false, true)

	return
}

// MessageParams ...
func (e *Actor) MessageParams() m.Attributes {
	return NewMessageParams()
}

func (e *Actor) Save(msg notifyCommon.Message) (addresses []string, message *m.Message) {
	message = &m.Message{
		Type:       Name,
		Attributes: msg.Attributes,
	}
	var err error
	if message.Id, err = e.Service.Adaptors().Message.Add(context.Background(), message); err != nil {
		log.Error(err.Error())
	}

	attr := NewMessageParams()
	_, _ = attr.Deserialize(message.Attributes)

	addresses = []string{fmt.Sprintf("%d", attr[AttrEntityID].Int64())}
	return
}

// Send ...
func (e *Actor) Send(address string, message *m.Message) (err error) {

	attr := NewMessageParams()
	if _, err = attr.Deserialize(message.Attributes); err != nil {
		log.Error(err.Error())
		return
	}

	chatID, _ := strconv.ParseInt(address, 0, 64)

	go func() {
		if _, _, err = e.client.SendMsg(attr[AttrBody].String(), chatID); err != nil {
			log.Error(err.Error())
		}
	}()
	return
}

func (e *Actor) eventHandler(topic string, msg interface{}) {

	switch v := msg.(type) {
	case notifyCommon.Message:
		if v.EntityId != nil && v.EntityId.PluginName() == Name {
			e.notify.SaveAndSend(v, e)
		}
	}
}
