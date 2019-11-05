package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type MessageDeliveries struct {
	Db *gorm.DB
}

type MessageDelivery struct {
	Id                 int64 `gorm:"primary_key"`
	Message            *Message
	MessageId          int64
	Address            string
	Status             string
	ErrorMessageStatus *string `gorm:"column:error_system_code"`
	ErrorMessageBody   *string `gorm:"column:error_system_message"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (d *MessageDelivery) TableName() string {
	return "message_deliveries"
}

func (n MessageDeliveries) Add(msg *MessageDelivery) (id int64, err error) {
	if err = n.Db.Create(&msg).Error; err != nil {
		return
	}
	id = msg.Id
	return
}

func (n *MessageDeliveries) List(limit, offset int64, orderBy, sort string) (list []*MessageDelivery, total int64, err error) {

	if err = n.Db.Model(MessageDelivery{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*MessageDelivery, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Preload("Message").
		Find(&list).
		Error

	return
}

func (n *MessageDeliveries) GetAllUncompleted(limit, offset int64) (list []*MessageDelivery, total int64, err error) {

	if err = n.Db.Model(MessageDelivery{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*MessageDelivery, 0)
	err = n.Db.
		Where("status in ('in_progress', 'new')").
		Limit(limit).
		Offset(offset).
		Preload("Message").
		Find(&list).
		Error

	return
}

func (n MessageDeliveries) SetStatus(msg *MessageDelivery) (err error) {

	err = n.Db.Model(&MessageDelivery{Id: msg.Id}).
		Updates(map[string]interface{}{
			"status":               msg.Status,
			"error_system_code":    msg.ErrorMessageStatus,
			"error_system_message": msg.ErrorMessageBody,
		}).Error
	return
}
