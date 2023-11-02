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

package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/atomic"
	tele "gopkg.in/telebot.v3"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/notify"
	"github.com/e154/smart-home/system/supervisor"
	"github.com/e154/smart-home/version"
)

// Actor ...
type Actor struct {
	supervisor.BaseActor
	isStarted   *atomic.Bool
	AccessToken string
	bot         *tele.Bot
	actionPool  chan events.EventCallEntityAction
	notify      *notify.Notify
}

// NewActor ...
func NewActor(entity *m.Entity,
	service supervisor.Service) (*Actor, error) {

	settings := NewSettings()
	_, _ = settings.Deserialize(entity.Settings.Serialize())

	actor := &Actor{
		BaseActor:   supervisor.NewBaseActor(entity, service),
		actionPool:  make(chan events.EventCallEntityAction, 1000),
		isStarted:   atomic.NewBool(false),
		AccessToken: settings[AttrToken].Decrypt(),
		notify:      notify.NewNotify(service.Adaptors()),
	}

	if actor.Attrs == nil {
		actor.Attrs = NewAttr()
	}

	if actor.Setts == nil {
		actor.Setts = NewSettings()
	}

	// Actions
	for _, a := range actor.Actions {
		if a.ScriptEngine.Engine() != nil {
			// bind
			_, _ = a.ScriptEngine.Engine().Do()
		}
	}

	// action worker
	go func() {
		for msg := range actor.actionPool {
			actor.runAction(msg)
		}
	}()

	return actor, nil
}

func (e *Actor) Destroy() {
	if !e.isStarted.Load() {
		return
	}
	e.Service.EventBus().Unsubscribe(notify.TopicNotify, e.eventHandler)
	e.notify.Shutdown()

	if e.bot != nil {
		e.bot.Stop()
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

	if !common.TestMode() {

		pref := tele.Settings{
			Token:  e.AccessToken,
			Poller: &tele.LongPoller{Timeout: 10 * time.Second},
		}

		e.bot, err = tele.NewBot(pref)
		if err != nil {
			err = errors.Wrap(apperr.ErrInternal, err.Error())
			return
		}

		e.bot.Handle("/help", e.commandHelp)
		e.bot.Handle("/start", e.commandStart)
		e.bot.Handle("/quit", e.commandQuit)
		e.bot.Handle(tele.OnText, e.commandAction)

		go e.bot.Start()
	}

	e.Service.EventBus().Subscribe(notify.TopicNotify, e.eventHandler, false)
	e.notify.Start()

	return
}

func (e *Actor) sendMsg(message *m.Message, chatId int64) (messageID int, err error) {

	var msg *tele.Message
	defer func() {
		if err == nil {
			if msg != nil {
				messageID = msg.ID
			}
			//go func() { _ = e.UpdateStatus() }()
			log.Infof("Sent message '%s' to chatId '%d'", message.Attributes, chatId)
		}
	}()
	if common.TestMode() {
		messageID = 123
		return
	}
	var chat *tele.Chat
	if chat, err = e.bot.ChatByID(chatId); err != nil {
		log.Error(err.Error())
		return
	}

	params := NewMessageParams()
	_, _ = params.Deserialize(message.Attributes)

	var body interface{}
	body = params[AttrBody].String()

	if uri := params[AttrUri].String(); uri != "" {
		body = &tele.Photo{File: tele.FromURL(uri)}
	}

	if uri := params[AttrFilePath].String(); uri != "" {
		body = &tele.Photo{File: tele.FromDisk(uri)}
	}

	msg, err = e.bot.Send(chat, body)
	return
}

func (e *Actor) getChatList() (list []m.TelegramChat, err error) {
	list, _, err = e.Service.Adaptors().TelegramChat.List(context.Background(), 999, 0, "", "", e.Id)
	return
}

// UpdateStatus ...
func (e *Actor) UpdateStatus() (err error) {

	oldState := e.GetEventState()
	now := e.Now(oldState)

	var attributeValues = make(m.AttributeValue)
	// ...

	e.AttrMu.Lock()
	var changed bool
	if changed, err = e.Attrs.Deserialize(attributeValues); !changed {
		if err != nil {
			log.Warn(err.Error())
		}

		if oldState.LastUpdated != nil {
			delta := now.Sub(*oldState.LastUpdated).Milliseconds()
			//fmt.Println("delta", delta)
			if delta < 200 {
				e.AttrMu.Unlock()
				return
			}
		}
	}
	e.AttrMu.Unlock()

	go e.SaveState(events.EventStateChanged{
		StorageSave: true,
		PluginName:  e.Id.PluginName(),
		EntityId:    e.Id,
		OldState:    oldState,
		NewState:    e.GetEventState(),
	})

	return
}

func (e *Actor) commandStart(c tele.Context) (err error) {

	var (
		user = c.Sender()
		chat = c.Chat()
		text = c.Text()
	)

	if pin := e.Setts[AttrPin].Decrypt(); pin != "" {
		if enterdPin := strings.Replace(text, "/start ", "", -1); pin != enterdPin {
			log.Warn("pin not equal")
			return
		}
	}

	text = fmt.Sprintf(banner, version.GetHumanVersion(), text)
	_ = e.Service.Adaptors().TelegramChat.Add(context.Background(), m.TelegramChat{
		EntityId: e.Id,
		ChatId:   chat.ID,
		Username: user.Username,
	})
	log.Infof("user '%s' added to chat", user.Username)
	err = c.Send(text, e.genKeyboard())
	return
}

func (e *Actor) commandHelp(c tele.Context) (err error) {

	builder := &strings.Builder{}
	if len(e.Actions) > 0 {
		for _, action := range e.Actions {
			builder.WriteString(fmt.Sprintf("/%s - %s\n", action.Name, action.Description))
		}
	}
	builder.WriteString(help)
	err = c.Send(builder.String(), e.genKeyboard())
	return err
}

func (e *Actor) commandQuit(c tele.Context) (err error) {

	var (
		chat = c.Chat()
	)

	_ = e.Service.Adaptors().TelegramChat.Delete(context.Background(), e.Id, chat.ID)
	err = c.Send("/quit -unsubscribe from bot\n/start - subscriber again")
	return
}

func (e *Actor) commandAction(c tele.Context) (err error) {

	var (
		text = c.Text()
	)

	e.runAction(events.EventCallEntityAction{
		ActionName: text,
		EntityId:   e.Id,
		Args: map[string]interface{}{
			"chatId":   c.Chat().ID,
			"username": c.Chat().Username,
		},
	})
	return
}

func (e *Actor) addAction(event events.EventCallEntityAction) {
	e.actionPool <- event
}

func (e *Actor) runAction(msg events.EventCallEntityAction) {
	actionName := strings.Replace(msg.ActionName, "/", "", 1)
	if action, ok := e.Actions[actionName]; ok {
		if action.ScriptEngine.Engine() == nil {
			return
		}
		if _, err := action.ScriptEngine.Engine().AssertFunction(FuncEntityAction, msg.EntityId, action.Name, msg.Args); err != nil {
			log.Error(err.Error())
			return
		}
		return
	}

	if e.ScriptEngines != nil {
		for _, engine := range e.ScriptEngines {
			if _, err := engine.Engine().AssertFunction(FuncEntityAction, msg.EntityId, msg.ActionName, msg.Args); err != nil {
				log.Error(err.Error())
				return
			}
		}
	}
}

// gen keyboard from actions
// [button][button][button]
// [button][button][button]
// [button][button][button]
func (e *Actor) genKeyboard() (menu *tele.ReplyMarkup) {
	menu = &tele.ReplyMarkup{ResizeKeyboard: true}
	var row []tele.Btn
	if len(e.Actions) == 0 {
		return
	}
	for k := range e.Actions {
		row = append(row, menu.Text(fmt.Sprintf("/%s", k)))
	}
	menu.Reply(menu.Split(3, row)...)
	return
}

// todo: prepare state
func (e *Actor) updateState(connected bool) {
	info := e.Info()
	var newStat = AttrOffline
	if connected {
		newStat = AttrConnected
	}
	if info.State != nil && info.State.Name == newStat {
		return
	}
	e.SetState(supervisor.EntityStateParams{
		NewState:    common.String(newStat),
		StorageSave: true,
	})
}

// Save ...
func (e *Actor) Save(msg notify.Message) (addresses []string, message *m.Message) {
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

	params := NewMessageParams()
	_, _ = params.Deserialize(message.Attributes)

	if val := params[AttrChatID].Int64(); val != 0 {
		addresses = []string{fmt.Sprintf("%d", val)}
	} else {
		addresses = []string{"broadcast"}
	}

	return
}

// Send ...
func (e *Actor) Send(address string, message *m.Message) (err error) {
	if !e.isStarted.Load() {
		return
	}

	var chatID *int64
	if address != "" && address != "broadcast" {
		var val int64
		if val, err = strconv.ParseInt(address, 10, 64); err == nil {
			chatID = common.Int64(val)
		}
	}

	if chatID != nil {
		if _, err = e.sendMsg(message, *chatID); err != nil {
			log.Warn(err.Error())
		}
		return
	}

	var list []m.TelegramChat
	if list, err = e.getChatList(); err != nil {
		return
	}
	for _, chat := range list {
		if _, err = e.sendMsg(message, chat.ChatId); err != nil {
			log.Warn(err.Error())
		}
	}

	return
}

// MessageParams ...
func (e *Actor) MessageParams() m.Attributes {
	return NewMessageParams()
}

func (e *Actor) eventHandler(topic string, msg interface{}) {

	switch v := msg.(type) {
	case notify.Message:
		if v.EntityId != nil && v.EntityId.PluginName() == Name {
			e.notify.SaveAndSend(v, e)
		}
	}
}
