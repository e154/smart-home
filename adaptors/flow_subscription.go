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

func (f *FlowSubscription) Add(sub *m.FlowSubscription) (err error) {
	err = f.table.Add(f.toDb(sub))
	return
}

func (f *FlowSubscription) Remove(ids []int64) (err error) {
	err = f.table.Delete(ids)
	return
}

func (f *FlowSubscription) fromDb(dbVer *db.FlowSubscription) (ver *m.FlowSubscription) {

	ver = &m.FlowSubscription{
		Id:     dbVer.Id,
		FlowId: dbVer.FlowId,
		Topic:  dbVer.Topic,
	}

	return
}

func (f *FlowSubscription) toDb(ver *m.FlowSubscription) (dbVer *db.FlowSubscription) {

	dbVer = &db.FlowSubscription{
		Id:     ver.Id,
		FlowId: ver.FlowId,
		Topic:  ver.Topic,
	}

	return
}
