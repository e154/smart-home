package adaptors

import (
	"github.com/jinzhu/gorm"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
)

type Image struct {
	table *db.Images
	db    *gorm.DB
}

func GetImageAdaptor(d *gorm.DB) *Image {
	return &Image{
		table: &db.Images{Db: d},
		db:    d,
	}
}

func (n *Image) fromDb(dbImage *db.Image) (image *m.Image) {
	image = &m.Image{
		Id:        dbImage.Id,
		Thumb:     dbImage.Thumb,
		Url:       dbImage.Url,
		Image:     dbImage.Image,
		MimeType:  dbImage.MimeType,
		Title:     dbImage.Title,
		Size:      dbImage.Size,
		Name:      dbImage.Name,
		CreatedAt: dbImage.CreatedAt,
	}
	return
}

func (n *Image) toDb(image *m.Image) (dbImage *db.Image) {
	dbImage = &db.Image{
		Id:        image.Id,
		Thumb:     image.Thumb,
		Url:       image.Url,
		Image:     image.Image,
		MimeType:  image.MimeType,
		Title:     image.Title,
		Size:      image.Size,
		Name:      image.Name,
		CreatedAt: image.CreatedAt,
	}
	return
}
