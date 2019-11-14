package notify

import "github.com/e154/smart-home/models"

type NotifyStat struct {
	MbBalance float32 `json:"mb_balance,omitempty"`
	TwBalance float32 `json:"tw_balance,omitempty"`
	Workers   int     `json:"workers"`
}

type IMessage interface {
	Save() (addresses []string, message *models.Message)
}
