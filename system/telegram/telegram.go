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

				if c.chatId == nil {
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

	if c.chatId != nil {
		msg := tgbotapi.NewMessage(common.Int64Value(c.chatId), body)
		_, err := c.bot.Send(msg)
		return err
	}

	return nil
}
