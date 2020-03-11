// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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
	"github.com/Syfaro/telegram-bot-api"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/logging"
	"github.com/pkg/errors"
)

var (
	log = common.MustGetLogger("telegram")
)

type Telegram struct {
	bot          *tgbotapi.BotAPI
	isStarted    bool
	stopPrecess  bool
	stopQueue    chan struct{}
	chatId       *int64
	updateChatId func(chatId int64)
	commandPool  chan Command
}

func NewTelegram(cfg *TelegramConfig, updateChatId func(chatId int64)) (*Telegram, error) {

	if cfg.Token == "" {
		return nil, errors.New("bad parameters")
	}

	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	log.Infof("Authorized on account %s", bot.Self.UserName)

	client := &Telegram{
		bot:          bot,
		stopQueue:    make(chan struct{}),
		updateChatId: updateChatId,
		chatId:       cfg.ChatId,
		commandPool:  make(chan Command),
	}

	go client.start()

	return client, nil
}

func (c *Telegram) start() {

	if c.isStarted {
		return
	}

	c.isStarted = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := c.bot.GetUpdatesChan(u)
	if err != nil {
		log.Error(err.Error())
		return
	}

	go func() {
		for {
			select {
			case update := <-updates:

				userName := update.Message.From.UserName
				chatID := update.Message.Chat.ID
				text := update.Message.Text

				if c.chatId == nil || common.Int64Value(c.chatId) == 0 {
					c.chatId = common.Int64(chatID)
					if c.updateChatId != nil {
						c.updateChatId(chatID)
					}
				}

				c.commandPool <- Command{
					ChatId:   chatID,
					Text:     text,
					UserName: userName,
				}

			case <-c.stopQueue:
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case v := <-c.commandPool:
				c.commandHandler(v)
			case <-c.stopQueue:
				return
			}
		}
	}()
}

func (c *Telegram) Stop() {
	c.stopPrecess = true
	c.isStarted = false

	c.bot.StopReceivingUpdates()
	c.stopQueue <- struct{}{}
	c.stopQueue <- struct{}{}

	c.stopPrecess = false
}

func (c *Telegram) SendMsg(body string) error {

	if !c.isStarted {
		return errors.New("bot not started")
	}

	if c.chatId != nil && common.Int64Value(c.chatId) != 0 {
		msg := tgbotapi.NewMessage(common.Int64Value(c.chatId), body)
		_, err := c.bot.Send(msg)
		return err
	}

	return nil
}
