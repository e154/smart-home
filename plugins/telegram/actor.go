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
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/atomic"
	tele "gopkg.in/telebot.v3"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
	"github.com/e154/smart-home/version"
)

// Actor ...
type Actor struct {
	supervisor.BaseActor
	isStarted   *atomic.Bool
	eventBus    bus.Bus
	adaptors    *adaptors.Adaptors
	AccessToken string
	bot         *tele.Bot
	msgPool     chan string
	actionPool  chan events.EventCallEntityAction
}

// NewActor ...
func NewActor(entity *m.Entity,
	visor supervisor.Supervisor,
	scriptService scripts.ScriptService,
	eventBus bus.Bus,
	adaptors *adaptors.Adaptors) (*Actor, error) {

	settings := NewSettings()
	_, _ = settings.Deserialize(entity.Settings.Serialize())

	actor := &Actor{
		BaseActor:   supervisor.NewBaseActor(entity, scriptService, adaptors),
		eventBus:    eventBus,
		actionPool:  make(chan events.EventCallEntityAction, 10),
		isStarted:   atomic.NewBool(false),
		adaptors:    adaptors,
		AccessToken: settings[AttrToken].String(),
		msgPool:     make(chan string, 99),
	}

	actor.Supervisor = visor

	if actor.Attrs == nil {
		actor.Attrs = NewAttr()
	}

	if actor.Setts == nil {
		actor.Setts = NewSettings()
	}

	// Actions
	for _, a := range actor.Actions {
		if a.ScriptEngine != nil {
			// bind
			_, _ = a.ScriptEngine.Do()
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
func (p *Actor) Spawn() supervisor.PluginActor {
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

		pref := tele.Settings{
			Token:  p.AccessToken,
			Poller: &tele.LongPoller{Timeout: 10 * time.Second},
		}

		p.bot, err = tele.NewBot(pref)
		if err != nil {
			err = errors.Wrap(apperr.ErrInternal, err.Error())
			return
		}

		p.bot.Handle("/help", p.commandHelp)
		p.bot.Handle("/start", p.commandStart)
		p.bot.Handle("/quit", p.commandQuit)
		p.bot.Handle(tele.OnText, p.commandAction)

		go p.bot.Start()
	}

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
		p.bot.Stop()
	}
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

func (p *Actor) sendMsg(body string, chatId int64) (messageID int, err error) {
	defer func() {
		if err == nil {
			go func() { _ = p.UpdateStatus() }()
			log.Infof("Sent message '%s' to chatId '%d'", body, chatId)
		}
	}()
	if common.TestMode() {
		messageID = 123
		return
	}
	var chat *tele.Chat
	if chat, err = p.bot.ChatByID(chatId); err != nil {
		log.Error(err.Error())
		return
	}
	var msg *tele.Message
	if msg, err = p.bot.Send(chat, body); err != nil {
		log.Error(err.Error())
		return
	}
	messageID = msg.ID
	return
}

func (p *Actor) getChatList() (list []m.TelegramChat, err error) {
	list, _, err = p.adaptors.TelegramChat.List(context.Background(), 999, 0, "", "", p.Id)
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

	p.eventBus.Publish("system/entities/"+p.Id.String(), events.EventStateChanged{
		StorageSave: true,
		PluginName:  p.Id.PluginName(),
		EntityId:    p.Id,
		OldState:    oldState,
		NewState:    p.GetEventState(p),
	})

	return
}

func (p *Actor) commandStart(c tele.Context) (err error) {

	var (
		user = c.Sender()
		chat = c.Chat()
		text = c.Text()
	)

	text = fmt.Sprintf(banner, version.GetHumanVersion(), text)
	_ = p.adaptors.TelegramChat.Add(context.Background(), m.TelegramChat{
		EntityId: p.Id,
		ChatId:   chat.ID,
		Username: user.Username,
	})
	err = c.Send(text, p.genKeyboard())
	return
}

func (p *Actor) commandHelp(c tele.Context) (err error) {

	builder := &strings.Builder{}
	if len(p.Actions) > 0 {
		for _, action := range p.Actions {
			builder.WriteString(fmt.Sprintf("/%s - %s\n", action.Name, action.Description))
		}
	}
	builder.WriteString(help)
	err = c.Send(builder.String(), p.genKeyboard())
	return err
}

func (p *Actor) commandQuit(c tele.Context) (err error) {

	var (
		chat = c.Chat()
	)

	_ = p.adaptors.TelegramChat.Delete(context.Background(), p.Id, chat.ID)
	err = c.Send("/quit -unsubscribe from bot\n/start - subscriber again")
	return
}

func (p *Actor) commandAction(c tele.Context) (err error) {

	var (
		text = c.Text()
	)

	p.runAction(events.EventCallEntityAction{
		ActionName: strings.Replace(text, "/", "", 1),
		EntityId:   p.Id,
	})
	return
}

func (p *Actor) addAction(event events.EventCallEntityAction) {
	p.actionPool <- event
}

func (p *Actor) runAction(msg events.EventCallEntityAction) {
	action, ok := p.Actions[msg.ActionName]
	if !ok {
		log.Warnf("action %s not found", msg.ActionName)
		return
	}
	if action.ScriptEngine == nil {
		return
	}
	if _, err := action.ScriptEngine.AssertFunction(FuncEntityAction, msg.EntityId, action.Name, msg.Args); err != nil {
		log.Error(err.Error())
		return
	}
}

// gen keyboard from actions
// [button][button][button]
// [button][button][button]
// [button][button][button]
func (p *Actor) genKeyboard() (menu *tele.ReplyMarkup) {
	menu = &tele.ReplyMarkup{ResizeKeyboard: true}
	var row []tele.Btn
	if len(p.Actions) == 0 {
		return
	}
	for k := range p.Actions {
		row = append(row, menu.Text(fmt.Sprintf("/%s", k)))
	}
	menu.Reply(menu.Split(3, row)...)
	return
}
