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

package telegram

import (
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/version"
	"go.uber.org/atomic"
	"sync"
)

// Actor ...
type Actor struct {
	entity_manager.BaseActor
	isStarted   *atomic.Bool
	eventBus    event_bus.EventBus
	adaptors    *adaptors.Adaptors
	AccessToken string
	bot         *tgbotapi.BotAPI
	commandPool chan Command
	msgPool     chan string
}

// NewActor ...
func NewActor(entity *m.Entity,
	entityManager entity_manager.EntityManager,
	eventBus event_bus.EventBus,
	adaptors *adaptors.Adaptors) (*Actor, error) {

	settings := NewSettings()
	settings.Deserialize(entity.Settings.Serialize())

	actor := &Actor{
		BaseActor: entity_manager.BaseActor{
			Id:         entity.Id,
			Name:       Name,
			EntityType: entity.Type,
			AttrMu:     &sync.RWMutex{},
			Attrs:      NewAttr(),
			Manager:    entityManager,
			SettingsMu: &sync.RWMutex{},
			Setts:      settings,
		},
		isStarted:   atomic.NewBool(false),
		eventBus:    eventBus,
		adaptors:    adaptors,
		AccessToken: settings[AttrToken].String(),
		commandPool: make(chan Command, 99),
		msgPool:     make(chan string, 99),
	}

	return actor, nil
}

func (p *Actor) Spawn() entity_manager.PluginActor {
	return p
}

func (p *Actor) Start() (err error) {

	if p.isStarted.Load() {
		return
	}
	defer func() {
		if err == nil {
			p.isStarted.Store(true)
		}
	}()

	if !common.TestMode() {
		p.bot, err = tgbotapi.NewBotAPI(p.AccessToken)
		if err != nil {
			err = fmt.Errorf("telegram error: %s", err.Error())
			return
		}

		log.Infof("Authorized on account %s", p.bot.Self.UserName)

		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		var updates tgbotapi.UpdatesChannel
		if updates, err = p.bot.GetUpdatesChan(u); err != nil {
			log.Error(err.Error())
			return
		}

		go func() {
			for update := range updates {
				p.commandPool <- Command{
					ChatId:   update.Message.Chat.ID,
					Text:     update.Message.Text,
					UserName: update.Message.From.UserName,
				}
			}
		}()
	}

	go func() {
		for v := range p.commandPool {
			p.commandHandler(v)
		}
	}()

	go func() {
		var list []m.TelegramChat
		var err error
		for msg := range p.msgPool {
			if list, err = p.getChatList(); err != nil {
				continue
			}
			for _, chat := range list {
				if _, err = p.sendMsg(msg, chat.ChatId); err != nil {
					log.Warn(err.Error())
				}
			}
		}
	}()

	return
}

func (p *Actor) Stop() {
	if !p.isStarted.Load() {
		return
	}
	if p.bot != nil {
		p.bot.StopReceivingUpdates()
	}
	close(p.commandPool)
	close(p.msgPool)
	p.isStarted.Store(false)
}

// Send ...
func (p *Actor) Send(message string) (err error) {
	if !p.isStarted.Load() {
		return
	}
	p.msgPool <- message
	return
}

// GetStatus ...
func (p *Actor) GetStatus(smsId string) (string, error) {

	return "", nil
}

func (p *Actor) commandHandler(cmd Command) {
	switch cmd.Text {
	case "/start":
		p.commandStart(cmd)
	case "/help":
		p.commandHelp(cmd)
	case "/quit":
		p.commandQuit(cmd)
	default:
		log.Infof("unknown command user(%s) chatId(%d) command(%s)", cmd.UserName, cmd.ChatId, cmd.Text)
	}
}

func (p *Actor) sendMsg(body string, chatId int64) (messageID int, err error) {
	defer func() {
		if err == nil {
			go p.UpdateStatus()
			log.Debugf("Sent message '%s' to chatId '%d'", body, chatId)
		}
	}()
	if common.TestMode() {
		messageID = 123
		return
	}
	msgCfg := tgbotapi.NewMessage(chatId, body)
	var msg tgbotapi.Message
	if msg, err = p.bot.Send(msgCfg); err != nil {
		return
	}
	messageID = msg.MessageID
	return
}

func (p *Actor) getChatList() (list []m.TelegramChat, err error) {
	list, _, err = p.adaptors.TelegramChat.List(999, 0, "", "", p.Id)
	return
}

// UpdateStatus ...
func (p *Actor) UpdateStatus() (err error) {

	oldState := p.GetEventState(p)
	now := p.Now(oldState)

	var attributeValues = make(m.AttributeValue)
	// ...

	p.AttrMu.Lock()
	var changed bool
	if changed, err = p.Attrs.Deserialize(attributeValues); !changed {
		if err != nil {
			log.Warn(err.Error())
		}

		if oldState.LastUpdated != nil {
			delta := now.Sub(*oldState.LastUpdated).Milliseconds()
			//fmt.Println("delta", delta)
			if delta < 200 {
				p.AttrMu.Unlock()
				return
			}
		}
	}
	p.AttrMu.Unlock()

	p.eventBus.Publish(event_bus.TopicEntities, event_bus.EventStateChanged{
		StorageSave: true,
		Type:        p.Id.Type(),
		EntityId:    p.Id,
		OldState:    oldState,
		NewState:    p.GetEventState(p),
	})

	return
}

func (p *Actor) commandStart(cmd Command) {
	text := fmt.Sprintf(banner, version.GetHumanVersion(), cmd.Text)
	p.adaptors.TelegramChat.Add(m.TelegramChat{
		EntityId: p.Id,
		ChatId:   cmd.ChatId,
		Username: cmd.UserName,
	})
	p.sendMsg(text, cmd.ChatId)
}

func (p *Actor) commandHelp(cmd Command) {
	p.sendMsg(help, cmd.ChatId)
}

func (p *Actor) commandQuit(cmd Command) {
	p.adaptors.TelegramChat.Delete(p.Id, cmd.ChatId)
	p.sendMsg("/quit -unsubscribe from bot\n/start - subscriber again", cmd.ChatId)
}
