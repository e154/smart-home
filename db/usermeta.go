package db

import "github.com/jinzhu/gorm"

type UserMetas struct {
	Db *gorm.DB
}

type UserMeta struct {
	Id     int64 `gorm:"primary_key"`
	UserId int64
	Key    string
	Value  string
}

func (m *UserMeta) TableName() string {
	return "user_metas"
}

func (m *UserMetas) UpdateOrCreate(meta *UserMeta) (id int64, err error) {

	err = m.Db.Update(&UserMeta{}).
		Where("user_id = ? and key = ?", meta.UserId, meta.Key).
		Updates(map[string]interface{}{"value": meta.Value}).
		Error

	if err != nil {
		err = m.Db.Create(&meta).Error
		id = meta.Id
	}

	return
}
