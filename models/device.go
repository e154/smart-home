package models

import "time"

type Device struct {
	Id          int64         `orm:"pk;auto;column(id)" json:"id"`
	Device      *Device       `orm:"rel(fk);null" json:"device"`
	Node        *Node         `orm:"rel(fk);null" json:"node"`
	Address     *int          `orm:"" json:"address"`
	Baud        int           `orm:"size(11)" json:"baud"`
	Sleep       int64         `orm:"size(32)" json:"sleep"`
	Description string        `orm:"size(254)" json:"description" valid:"MaxSize(254)"`
	Name        string        `orm:"size(254)" json:"name" valid:"MaxSize(254);Required"`
	Status      string        `orm:"size(254)" json:"status" valid:"MaxSize(254)"`
	StopBite    int64         `orm:"size(11)" json:"stop_bite"`
	Timeout     time.Duration `orm:"" json:"timeout"`
	Tty         string        `orm:"size(254)" json:"tty" valid:"MaxSize(254)"`
	//States      []*DeviceState  `orm:"reverse(many)" json:"states"`
	//Actions     []*DeviceAction `orm:"reverse(many)" json:"actions"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime);column(created_at)" json:"created_at"`
	UpdateAt  time.Time `orm:"auto_now;type(datetime);column(update_at)" json:"update_at"`
	IsGroup   bool      `gorm:"-" json:"is_group"`
}
