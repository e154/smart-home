package db

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Messages struct {
	Db *gorm.DB
}

type Message struct {
	Id           int64 `gorm:"primary_key"`
	Type         string
	EmailFrom    *string
	EmailSubject *string
	EmailBody    *string
	SmsText      *string
	UiText       *string
	TelegramText *string
	Statuses     []*MessageDelivery
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (d *Message) TableName() string {
	return "messages"
}

func (n Messages) Add(msg *Message) (id int64, err error) {
	if err = n.Db.Create(&msg).Error; err != nil {
		return
	}
	id = msg.Id
	return
}
