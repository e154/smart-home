package adaptors

import (
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/gorm"
)

type FlowSubscription struct {
	db    *gorm.DB
	table *db.FlowSubscriptions
}

func GetFlowSubscriptionAdaptor(Db *gorm.DB) *FlowSubscription {
	return &FlowSubscription{
		db:    Db,
		table: db.NewFlowSubscriptions(Db),
	}
}

func (f *FlowSubscription) fromDb(dbVer *db.FlowSubscription) (ver *m.FlowSubscription) {

	ver = &m.FlowSubscription{
		Id:    dbVer.Id,
		Topic: dbVer.Topic,
	}

	return
}

func (f *FlowSubscription) toDb(ver *m.FlowSubscription) (dbVer *db.FlowSubscription) {

	dbVer = &db.FlowSubscription{
		Id:    ver.Id,
		Topic: ver.Topic,
	}

	return
}
