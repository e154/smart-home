package db

import (
	"github.com/jinzhu/gorm"
	"time"
)

type FlowSubscriptions struct {
	db *gorm.DB
}

func NewFlowSubscriptions(db *gorm.DB) *FlowSubscriptions {
	return &FlowSubscriptions{db: db}
}

type FlowSubscription struct {
	Id        int64 `gorm:"primary_key"`
	Flow      *Flow
	FlowId    int64
	Topic     string
	CreatedAt time.Time
}

func (d *FlowSubscription) TableName() string {
	return "flow_subscriptions"
}

func (f *FlowSubscriptions) Add(sub *FlowSubscription) (err error) {
	err = f.db.Create(sub).Error
	return
}

func (f *FlowSubscriptions) Delete(ids []int64) (err error) {
	err = f.db.Delete(&FlowSubscription{}, "id in (?)", ids).Error
	return
}
