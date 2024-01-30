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
	"path/filepath"
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
	notifyCommon "github.com/e154/smart-home/plugins/notify/common"
	"github.com/e154/smart-home/system/supervisor"
)

// Actor ...
type Actor struct {
	*supervisor.BaseActor
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
	_ = e.Service.EventBus().Unsubscribe(notify.TopicNotify, e.eventHandler)
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

		e.bot.Handle("/start", e.commandStart)
		e.bot.Handle("/quit", e.commandQuit)
		e.bot.Handle(tele.OnText, e.commandAction)

		go e.bot.Start()
	}

	_ = e.Service.EventBus().Subscribe(notify.TopicNotify, e.eventHandler, false)
	e.notify.Start()

	e.BaseActor.Spawn()
}

func (e *Actor) sendMsg(message *m.Message, chatId int64) (messageID int, err error) {

	var msg *tele.Message
	defer func() {
		if err == nil {
			if msg != nil {
				messageID = msg.ID
			}
			//go func() { _ = e.UpdateStatus() }()
			log.Infof("Sent message '%v' to chatId '%d'", message.Attributes, chatId)
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
	if _, err = params.Deserialize(message.Attributes); err != nil {
		return
	}

	var body interface{}

	keys := params[AttrKeys].ArrayString()

	// photos
	urls := params[AttrPhotoUri].ArrayString()
	if len(urls) > 0 {
		for _, uri := range urls {
			log.Infof("send photo %s", uri)
			if msg, err = e.bot.Send(chat, &tele.Photo{File: tele.FromURL(uri)}); err != nil {
				return
			}
		}
	}
	path := params[AttrPhotoPath].ArrayString()
	if len(path) > 0 {
		for _, uri := range path {
			log.Infof("send photo %s", uri)
			if msg, err = e.bot.Send(chat, &tele.Photo{File: tele.FromDisk(uri)}); err != nil {
				return
			}
		}
	}

	// files
	urls = params[AttrFileUri].ArrayString()
	if len(urls) > 0 {
		for _, uri := range urls {
			log.Infof("send file %s", uri)
			fileName := filepath.Base(uri)
			if msg, err = e.bot.Send(chat, &tele.Document{File: tele.FromURL(uri), FileName: fileName}); err != nil {
				return
			}
		}
	}
	path = params[AttrFilePath].ArrayString()
	if len(path) > 0 {
		for _, uri := range path {
			log.Infof("send file %s", uri)
			fileName := filepath.Base(uri)
			if msg, err = e.bot.Send(chat, &tele.Document{File: tele.FromDisk(uri), FileName: fileName}); err != nil {
				return
			}
		}
	}

	if body = params[AttrBody].String(); body != "" {
		msg, err = e.bot.Send(chat, body, e.genPlainKeyboard(keys))
	}
	return
}

func (e *Actor) getChatList() (list []m.TelegramChat, err error) {
	list, _, err = e.Service.Adaptors().TelegramChat.List(context.Background(), 999, 0, "", "", e.Id)
	return
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

func (e *Actor) commandStart(c tele.Context) (err error) {

	var (
		user = c.Sender()
		chat = c.Chat()
		text = c.Text()
	)

	if pin := e.Setts[AttrPin].Decrypt(); pin != "" {
		enterdPin := strings.Replace(text, "/start ", "", -1)
		if pin != enterdPin {
			log.Warnf("received start command with bad pin code: \"%s\", username \"%s\"", enterdPin, chat.Username)
			return
		}
	}

	_ = e.Service.Adaptors().TelegramChat.Add(context.Background(), m.TelegramChat{
		EntityId: e.Id,
		ChatId:   chat.ID,
		Username: user.Username,
	})
	log.Infof("user '%s' added to chat", user.Username)

	e.runAction(events.EventCallEntityAction{
		ActionName: "/start",
		EntityId:   e.Id,
		Args: map[string]interface{}{
			"chatId":    c.Chat().ID,
			"username":  c.Chat().Username,
			"firstName": c.Chat().FirstName,
			"lastName":  c.Chat().LastName,
		},
	})

	return
}

func (e *Actor) commandQuit(c tele.Context) (err error) {

	var (
		chat = c.Chat()
	)

	_ = e.Service.Adaptors().TelegramChat.Delete(context.Background(), e.Id, chat.ID)
	menu := &tele.ReplyMarkup{RemoveKeyboard: true}
	var message = "/start - subscriber again"
	if pin := e.Setts[AttrPin].Decrypt(); pin != "" {
		message = "/start [pin] - subscriber again"
	}
	err = c.Send(message, menu)
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
	if action, ok := e.Actions[msg.ActionName]; ok {
		if action.ScriptEngine != nil && action.ScriptEngine.Engine() != nil {
			if _, err := action.ScriptEngine.Engine().AssertFunction(FuncEntityAction, msg.EntityId, action.Name, msg.Args); err != nil {
				log.Error(errors.Wrapf(err, "entity id: %s ", e.Id).Error())
			}
			return
		}
	}
	if e.ScriptsEngine != nil && e.ScriptsEngine.Engine() != nil {
		if _, err := e.ScriptsEngine.AssertFunction(FuncEntityAction, msg.EntityId, msg.ActionName, msg.Args); err != nil {
			log.Error(errors.Wrapf(err, "entity id: %s ", e.Id).Error())
		}
	}
}

// gen keyboard from actions
// [button][button][button]
// [button][button][button]
// [button][button][button]
func (e *Actor) genActionKeyboard() (menu *tele.ReplyMarkup) {
	menu = &tele.ReplyMarkup{
		ResizeKeyboard: true,
		RemoveKeyboard: len(e.Actions) == 0,
	}
	var row []tele.Btn
	if len(e.Actions) == 0 {
		return
	}
	for _, action := range e.Actions {
		row = append(row, menu.Text(action.Name))
	}
	menu.Reply(menu.Split(3, row)...)
	return
}

func (e *Actor) genPlainKeyboard(keys []string) (menu *tele.ReplyMarkup) {
	menu = &tele.ReplyMarkup{
		ResizeKeyboard: true,
		RemoveKeyboard: len(keys) == 0,
	}
	var row []tele.Btn
	for _, key := range keys {
		row = append(row, menu.Text(key))
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
	_ = e.SetState(supervisor.EntityStateParams{
		NewState:    common.String(newStat),
		StorageSave: true,
	})
}

// Save ...
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
	case notifyCommon.Message:
		if v.EntityId != nil && v.EntityId.PluginName() == Name {
			e.notify.SaveAndSend(v, e)
		}
	}
}
