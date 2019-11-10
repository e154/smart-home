package telegram

import (
	"github.com/Syfaro/telegram-bot-api"
	"github.com/e154/smart-home/common"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
)

var (
	log = logging.MustGetLogger("telegram")
)

type Telegram struct {
	bot          *tgbotapi.BotAPI
	isStarted    bool
	stopPrecess  bool
	stopQueue    chan struct{}
	chatId       *int64
	updateChatId func(chatId int64)
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

	for {
		select {
		case update := <-updates:

			// Пользователь, который написал боту
			UserName := update.Message.From.UserName

			// ID чата/диалога.
			// Может быть идентификатором как чата с пользователем
			// (тогда он равен UserID) так и публичного чата/канала
			chatID := update.Message.Chat.ID

			if c.chatId == nil {
				c.chatId = common.Int64(chatID)
				if c.updateChatId != nil {
					c.updateChatId(chatID)
				}
			}

			// Текст сообщения
			Text := update.Message.Text

			log.Infof("[%s] %d %s", UserName, chatID, Text)

			// Ответим пользователю его же сообщением
			reply := Text

			// Созадаем сообщение
			msg := tgbotapi.NewMessage(chatID, reply)

			// и отправляем его
			c.bot.Send(msg)

		case <-c.stopQueue:
			return

		}
	}
}

func (c *Telegram) Stop() {
	c.stopPrecess = true
	c.isStarted = false

	c.bot.StopReceivingUpdates()
	c.stopQueue <- struct{}{}

	c.stopPrecess = false
}

func (c *Telegram) SendMsg(body string, channels []string) error {

	if !c.isStarted {
		return errors.New("bot not started")
	}

	return nil
}
