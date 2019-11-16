package telegram

import (
	"fmt"
	"github.com/e154/smart-home/version"
)

const banner = `
Smart home system

Version:
%s

command:
%s
`

func (c *Telegram) commandHandler(cmd Command) {
	switch cmd.Text {
	case "/start":
		c.commandStart(cmd)
	case "/help":
		c.commandHelp(cmd)
	default:
		log.Infof("[%s] %d %s", cmd.UserName, cmd.ChatId, cmd.Text)
	}
}

func (c *Telegram) commandStart(cmd Command) {

	text := fmt.Sprintf(banner, version.GetHumanVersion(), cmd.Text)
	c.SendMsg(text)
}

func (c *Telegram) commandHelp(cmd Command) {

	text := fmt.Sprintf(banner, version.GetHumanVersion(), cmd.Text)
	c.SendMsg(text)
}
