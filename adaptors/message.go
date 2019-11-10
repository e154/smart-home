package adaptors

import (
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/gorm"
)

type Message struct {
	table *db.Messages
}

func GetMessageAdaptor(d *gorm.DB) *Message {
	return &Message{
		table: &db.Messages{Db: d},
	}
}

func (n *Message) Add(msg *m.Message) (id int64, err error) {
	id, err = n.table.Add(n.toDb(msg))
	return
}

func (n *Message) fromDb(dbVer *db.Message) (ver *m.Message) {
	ver = &m.Message{
		Id:           dbVer.Id,
		Type:         m.MessageType(dbVer.Type),
		EmailFrom:    dbVer.EmailFrom,
		EmailSubject: dbVer.EmailSubject,
		EmailBody:    dbVer.EmailBody,
		SmsText:      dbVer.SmsText,
		SlackText:    dbVer.SlackText,
		UiText:       dbVer.UiText,
		TelegramText: dbVer.TelegramText,
		CreatedAt:    dbVer.CreatedAt,
		UpdatedAt:    dbVer.UpdatedAt,
	}
	return
}

func (n *Message) toDb(ver *m.Message) (dbVer *db.Message) {
	dbVer = &db.Message{
		Id:           ver.Id,
		Type:         string(ver.Type),
		EmailFrom:    ver.EmailFrom,
		EmailSubject: ver.EmailSubject,
		EmailBody:    ver.EmailBody,
		SmsText:      ver.SmsText,
		SlackText:    ver.SlackText,
		UiText:       ver.UiText,
		TelegramText: ver.TelegramText,
		CreatedAt:    ver.CreatedAt,
		UpdatedAt:    ver.UpdatedAt,
	}
	return
}
