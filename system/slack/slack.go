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

package slack

import (
	"github.com/e154/smart-home/common"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
)

var (
	log = common.MustGetLogger("slack")
)

// Slack ...
type Slack struct {
	cfg *SlackConfig
	api *slack.Client
}

// NewSlack ...
func NewSlack(cfg *SlackConfig) (*Slack, error) {

	if cfg.Token == "" {
		return nil, errors.New("bad parameters")
	}

	return &Slack{
		cfg: cfg,
		api: slack.New(cfg.Token),
	}, nil
}

// SendMsg ...
func (c *Slack) SendMsg(message *SlackMessage) (err error) {

	options := []slack.MsgOption{
		slack.MsgOptionText(message.Text, false),
	}

	if c.cfg.UserName != "" {
		options = append(options, slack.MsgOptionUsername(c.cfg.UserName))
	}

	var channelID, timestamp string
	if channelID, timestamp, err = c.api.PostMessage(message.Channel, options...); err != nil {
		log.Error(err.Error())
		return
	}
	log.Infof("Message successfully sent to channel %s at %s", channelID, timestamp)
	return
}
