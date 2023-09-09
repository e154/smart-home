// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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

package webpush

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/e154/smart-home/system/supervisor"

	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/notify"
)

var (
	log = logger.MustGetLogger("plugins.webpush")
)

var _ supervisor.Pluggable = (*plugin)(nil)

func init() {
	supervisor.RegisterPlugin(Name, New)
}

type plugin struct {
	*supervisor.Plugin
	VAPIDPublicKey, VAPIDPrivateKey string
}

// New ...
func New() supervisor.Pluggable {
	return &plugin{
		Plugin: supervisor.NewPlugin(),
	}
}

// Load ...
func (p *plugin) Load(ctx context.Context, service supervisor.Service) (err error) {
	if err = p.Plugin.Load(ctx, service); err != nil {
		return
	}

	go func() {
		if err = p.asyncLoad(); err != nil {
			log.Error(err.Error())
		}
	}()

	return nil
}

func (p *plugin) asyncLoad() (err error) {

	// load settings
	var settings m.Attributes
	settings, err = p.LoadSettings(p)
	if err != nil {
		log.Warn(err.Error())
		settings = NewSettings()
	}

	if settings == nil {
		settings = NewSettings()
	}
	if settings[AttrPrivateKey].String() == "" || settings[AttrPublicKey].String() == "" {
		log.Info(`Keys not found, will be generate`)

		if settings[AttrPrivateKey].Value, settings[AttrPublicKey].Value, err = GenerateVAPIDKeys(); err != nil {
			return
		}
		var model *m.Plugin
		model, _ = p.Adaptors.Plugin.GetByName(context.Background(), Name)
		model.Settings = settings.Serialize()
		_ = p.Adaptors.Plugin.Update(context.Background(), model)
	}

	p.VAPIDPrivateKey = settings[AttrPrivateKey].String()
	p.VAPIDPublicKey = settings[AttrPublicKey].String()

	log.Infof(`Used public key: "%s"`, p.VAPIDPublicKey)

	_ = p.EventBus.Subscribe(TopicPluginWebpush, p.eventHandler)

	// register webpush provider
	notify.ProviderManager.AddProvider(Name, p)

	return
}

// Unload ...
func (p *plugin) Unload(ctx context.Context) (err error) {
	if err = p.Plugin.Unload(ctx); err != nil {
		return
	}

	notify.ProviderManager.RemoveProvider(Name)

	_ = p.EventBus.Unsubscribe(TopicPluginWebpush, p.eventHandler)

	return nil
}

// Name ...
func (p *plugin) Name() string {
	return Name
}

// Type ...
func (p *plugin) Type() supervisor.PluginType {
	return supervisor.PluginInstallable
}

// Depends ...
func (p *plugin) Depends() []string {
	return []string{notify.Name}
}

// Version ...
func (p *plugin) Version() string {
	return "0.0.1"
}

// Options ...
func (p *plugin) Options() m.PluginOptions {
	return m.PluginOptions{
		Setts: NewSettings(),
	}
}

// Save ...
func (p *plugin) Save(msg notify.Message) (addresses []string, message *m.Message) {
	message = &m.Message{
		Type:       Name,
		Attributes: msg.Attributes,
	}
	var err error
	if message.Id, err = p.Adaptors.Message.Add(context.Background(), message); err != nil {
		log.Error(err.Error())
	}

	attr := NewMessageParams()
	_, _ = attr.Deserialize(message.Attributes)

	addresses = strings.Split(attr[AttrUserIDS].String(), ",")

	return
}

// Send ...
func (p *plugin) Send(address string, message *m.Message) (err error) {

	attr := NewMessageParams()
	if _, err = attr.Deserialize(message.Attributes); err != nil {
		log.Error(err.Error())
		return
	}

	userId, _ := strconv.ParseInt(address, 0, 64)
	var userDevices []*m.UserDevice
	if userDevices, err = p.Adaptors.UserDevice.GetByUserId(userId); err != nil {
		return
	}

	go func() {
		for _, device := range userDevices {
			if err = p.sendPush(device, attr[AttrTitle].String(), attr[AttrBody].String()); err != nil {
				log.Error(err.Error())
			}
		}
	}()

	return
}

// MessageParams ...
func (p *plugin) MessageParams() m.Attributes {
	return NewMessageParams()
}

func (p *plugin) sendPush(userDevice *m.UserDevice, msgTitle, msgBody string) (err error) {

	msg := map[string]string{
		"title": msgTitle,
		"body":  msgBody,
	}

	message, _ := json.Marshal(msg)

	var statusCode int
	var responseBody []byte
	statusCode, responseBody, err = SendNotification(message, userDevice.Subscription, &Options{
		Crawler:         p.Crawler,
		VAPIDPublicKey:  p.VAPIDPublicKey,
		VAPIDPrivateKey: p.VAPIDPrivateKey,
		TTL:             30,
	})
	if err != nil {
		return
	}

	if statusCode != 201 {
		log.Warn(string(responseBody))
		go func() {
			_ = p.Adaptors.UserDevice.Delete(userDevice.Id)
			log.Infof("remove user device %d", userDevice.Id)
		}()
		return
	}

	log.Infof(`Send push, user: "%d", device: "%d", title: "%s"`, userDevice.UserId, userDevice.Id, msgTitle)

	return
}

func (p *plugin) eventHandler(_ string, event interface{}) {

	switch v := event.(type) {
	case EventAddWebPushSubscription:
		p.updateSubscribe(v)
	case EventGetWebPushPublicKey:
		p.sendPublicKey(v)
	}
}

func (p *plugin) sendPublicKey(event EventGetWebPushPublicKey) {
	p.EventBus.Publish(TopicPluginWebpush, EventNewWebPushPublicKey{
		UserID:    event.UserID,
		PublicKey: p.VAPIDPublicKey,
	})
}

func (p *plugin) updateSubscribe(event EventAddWebPushSubscription) {

	if _, err := p.Adaptors.UserDevice.Add(&m.UserDevice{
		UserId:       event.UserID,
		Subscription: event.Subscription,
	}); err != nil {
		return
	}

	log.Infof("user subscription updated, %d", event.UserID)
}
