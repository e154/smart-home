package db

import "time"

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
