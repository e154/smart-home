package db

import (
	"time"
	"github.com/jinzhu/gorm"
)

type Images struct {
	Db *gorm.DB
}

type Image struct {
	Id        int64 `gorm:"primary_key"`
	Thumb     string
	Url       string `gorm:"-"`
	Image     string
	MimeType  string
	Title     string
	Size      int64
	Name      string
	CreatedAt time.Time
}

func (m *Image) TableName() string {
	return "images"
}