package slack

import (
	"github.com/nlopes/slack"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
)

var (
	log = logging.MustGetLogger("slack")
)

type Slack struct {
	cfg *SlackConfig
	api *slack.Client
}

func NewSlack(cfg *SlackConfig) (*Slack, error) {

	if cfg.Token == "" {
		return nil, errors.New("bad parameters")
	}

	return &Slack{
		cfg: cfg,
		api: slack.New(cfg.Token),
	}, nil
}

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
