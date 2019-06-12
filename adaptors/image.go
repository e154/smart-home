package adaptors

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"github.com/t-tiger/gorm-bulk-insert"
	"github.com/jinzhu/gorm"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"strings"
	"path"
	"fmt"
	"os"
	"path/filepath"
	"github.com/e154/smart-home/common"
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

func (n *Image) Add(ver *m.Image) (id int64, err error) {

	dbVer := n.toDb(ver)
	if id, err = n.table.Add(dbVer); err != nil {
		return
	}

	return
}

func (n *Image) GetById(mapId int64) (ver *m.Image, err error) {

	var dbVer *db.Image
	if dbVer, err = n.table.GetById(mapId); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

func (n *Image) Update(ver *m.Image) (err error) {
	dbVer := n.toDb(ver)
	err = n.table.Update(dbVer)
	return
}

func (n *Image) Delete(mapId int64) (err error) {
	err = n.table.Delete(mapId)
	return
}

func (n *Image) List(limit, offset int64, orderBy, sort string) (list []*m.Image, total int64, err error) {
	var dbList []*db.Image
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.Image, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

func (n *Image) UploadImage(reader *bufio.Reader, fileName string) (err error) {

	buffer := bytes.NewBuffer(make([]byte, 0))
	part := make([]byte, 128)

	var count int
	for {
		if count, err = reader.Read(part); err != nil {
			break
		}
		buffer.Write(part[:count])
	}
	if err != io.EOF {
	} else {
		err = nil
	}

	contentType := http.DetectContentType(buffer.Bytes())
	log.Infof("Content-type from buffer, %s", contentType)

	//------
	// rename & save
	name := common.Strtomd5(common.RandomString(10))
	ext := strings.ToLower(path.Ext(fileName))
	newname := fmt.Sprintf("%s%s", name, ext)

	//create destination file making sure the path is writeable.
	dir := common.GetFullPath(name)
	os.MkdirAll(dir, os.ModePerm)
	var dst *os.File
	if dst, err = os.Create(filepath.Join(dir, newname)); err != nil {
		return
	}

	defer dst.Close()

	//copy the uploaded file to the destination file
	if _, err = io.Copy(dst, buffer); err != nil {
		return
	}

	size, _ := common.GetFileSize(filepath.Join(dir, newname))
	newFile := &m.Image{
		Size:     size,
		MimeType: contentType,
		Image:    newname,
		Name:     fileName,
	}

	newFile.Id, err = n.Add(newFile)

	return
}

func (n *Image) AddMultiple(items []*m.Image) (err error) {

	insertRecords := make([]interface{}, 0)
	for _, ver := range items {
		dbVer := n.toDb(ver)
		insertRecords = append(insertRecords, dbVer)
	}

	err = gormbulk.BulkInsert(n.db, insertRecords, 3000)

	return
}

func (n *Image) GetAllByDate(filter string) (images []*m.Image, err error) {

	var dblist []*db.Image
	if dblist, err = n.table.GetAllByDate(filter); err != nil {
		return
	}
	for _, dbVer := range dblist {
		ver := n.fromDb(dbVer)
		images = append(images, ver)
	}

	return
}

func (n *Image) GetFilterList() (filterList []*m.ImageFilterList, err error) {

	var dblist []*db.ImageFilterList
	if dblist, err = n.table.GetFilterList(); err != nil {
		return
	}
	for _, dbVer := range dblist {
		ver := &m.ImageFilterList{
			Date:  dbVer.Date,
			Count: dbVer.Count,
		}
		filterList = append(filterList, ver)
	}
	return
}

func (n *Image) fromDb(dbImage *db.Image) (image *m.Image) {
	image = &m.Image{
		Id:        dbImage.Id,
		Thumb:     dbImage.Thumb,
		Image:     dbImage.Image,
		MimeType:  dbImage.MimeType,
		Title:     dbImage.Title,
		Size:      dbImage.Size,
		Name:      dbImage.Name,
		CreatedAt: dbImage.CreatedAt,
	}
	if image.Image != "" {
		image.Url = common.GetLinkPath(image.Image)
	}
	return
}

func (n *Image) toDb(image *m.Image) (dbImage *db.Image) {
	dbImage = &db.Image{
		Id:       image.Id,
		Thumb:    image.Thumb,
		Image:    image.Image,
		MimeType: image.MimeType,
		Title:    image.Title,
		Size:     image.Size,
		Name:     image.Name,
	}
	return
}
