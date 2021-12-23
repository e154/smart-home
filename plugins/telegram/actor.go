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
	"strings"

	"github.com/pkg/errors"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/version"
	"go.uber.org/atomic"
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
	actionPool  chan event_bus.EventCallAction
}

// NewActor ...
func NewActor(entity *m.Entity,
	entityManager entity_manager.EntityManager,
	scriptService scripts.ScriptService,
	eventBus event_bus.EventBus,
	adaptors *adaptors.Adaptors) (*Actor, error) {

	settings := NewSettings()
	settings.Deserialize(entity.Settings.Serialize())

	actor := &Actor{
		BaseActor:   entity_manager.NewBaseActor(entity, scriptService, adaptors),
		eventBus:    eventBus,
		actionPool:  make(chan event_bus.EventCallAction, 10),
		isStarted:   atomic.NewBool(false),
		adaptors:    adaptors,
		AccessToken: settings[AttrToken].String(),
		commandPool: make(chan Command, 99),
		msgPool:     make(chan string, 99),
	}

	actor.Manager = entityManager

	if actor.Attrs == nil {
		actor.Attrs = NewAttr()
	}

	if actor.Setts == nil {
		actor.Setts = NewSettings()
	}

	//if actor.Actions == nil {
	//	actor.Actions = NewActions()
	//}

	// Actions
	for _, a := range actor.Actions {
		if a.ScriptEngine != nil {
			// bind
			a.ScriptEngine.Do()
		}
	}

	actor.DeserializeAttr(entity.Attributes.Serialize())

	// action worker
	go func() {
		for msg := range actor.actionPool {
			actor.runAction(msg)
		}
	}()

	return actor, nil
}

// Spawn ...
func (p *Actor) Spawn() entity_manager.PluginActor {
	return p
}

// Start ...
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
			err = errors.Wrap(common.ErrInternal, err.Error())
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

// Stop ...
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

func (p *Actor) commandHandler(cmd Command) {
	switch cmd.Text {
	case "/start":
		p.commandStart(cmd)
	case "/help":
		p.commandHelp(cmd)
	case "/quit":
		p.commandQuit(cmd)
	default:
		p.commandAction(cmd)
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
	p.genKeyboard(msgCfg)
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
		PluginName:  p.Id.PluginName(),
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
	builder := &strings.Builder{}
	if len(p.Actions) > 0 {
		for _, action := range p.Actions {
			builder.WriteString(fmt.Sprintf("/%s - %s\n", action.Name, action.Description))
		}
	}
	builder.WriteString(help)
	p.sendMsg(builder.String(), cmd.ChatId)
}

func (p *Actor) commandQuit(cmd Command) {
	p.adaptors.TelegramChat.Delete(p.Id, cmd.ChatId)
	p.sendMsg("/quit -unsubscribe from bot\n/start - subscriber again", cmd.ChatId)
}

//todo add command args
func (p *Actor) commandAction(cmd Command) {
	p.runAction(event_bus.EventCallAction{
		ActionName: strings.Replace(cmd.Text, "/", "", 1),
		EntityId:   p.Id,
	})
}

func (p *Actor) addAction(event event_bus.EventCallAction) {
	p.actionPool <- event
}

func (p *Actor) runAction(msg event_bus.EventCallAction) {
	action, ok := p.Actions[msg.ActionName]
	if !ok {
		log.Warnf("action %s not found", msg.ActionName)
		return
	}
	if _, err := action.ScriptEngine.AssertFunction(FuncEntityAction, msg.EntityId, action.Name); err != nil {
		log.Error(err.Error())
		return
	}
}

// gen keyboard from actions
// [button][button][button]
// [button][button][button]
// [button][button][button]
func (p *Actor) genKeyboard(msgCfg tgbotapi.MessageConfig) {
	var row []tgbotapi.KeyboardButton
	var rows [][]tgbotapi.KeyboardButton
	var counter = 0
	if len(p.Actions) == 0 {
		return
	}
	for k := range p.Actions {
		counter++
		if counter >= 3 {
			counter = 1
			rows = append(rows, row)
			row = []tgbotapi.KeyboardButton{}
		}
		row = append(row, tgbotapi.NewKeyboardButton(fmt.Sprintf("/%s", k)))
	}
	if counter < 3 {
		rows = append(rows, row)
	}
	msgCfg.ReplyMarkup = tgbotapi.NewReplyKeyboard(rows...)
}
