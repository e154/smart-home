package use_case

import (
	"github.com/e154/smart-home/system/validation"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/copier"
	"errors"
	"mime/multipart"
	"bufio"
)

type ImageCommand struct {
	*CommonCommand
}

func NewImageCommand(common *CommonCommand) *ImageCommand {
	return &ImageCommand{
		CommonCommand: common,
	}
}

func (i *ImageCommand) Add(params *m.Image) (image *m.Image, errs []*validation.Error, err error) {

	_, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	var id int64
	if id, err = i.adaptors.Image.Add(params); err != nil {
		return
	}

	image, err = i.adaptors.Image.GetById(id)

	return
}

func (i *ImageCommand) GetById(id int64) (image *m.Image, err error) {

	image, err = i.adaptors.Image.GetById(id)

	return
}

func (i *ImageCommand) Update(params *m.Image) (result *m.Image, errs []*validation.Error, err error) {

	var image *m.Image
	if image, err = i.adaptors.Image.GetById(params.Id); err != nil {
		return
	}

	if err = copier.Copy(&image, &params); err != nil {
		return
	}

	_, errs = image.Valid()
	if len(errs) > 0 {
		return
	}

	if err = i.adaptors.Image.Update(image); err != nil {
		return
	}

	image, err = i.adaptors.Image.GetById(params.Id)

	return
}

func (i *ImageCommand) Delete(imageId int64) (err error) {

	if imageId == 0 {
		err = errors.New("image id is null")
		return
	}

	var image *m.Image
	if image, err = i.adaptors.Image.GetById(imageId); err != nil {
		return
	}

	err = i.adaptors.Image.Delete(image.Id)

	return
}

func (i *ImageCommand) Upload(files map[string][]*multipart.FileHeader) (fileList []*m.Image, errs []error) {

	fileList = make([]*m.Image, 0)
	errs = make([]error, 0)

	for _, fileHeader := range files {

		file, err := fileHeader[0].Open()
		if err != nil {
			errs = append(errs, err)
			continue
		}

		reader := bufio.NewReader(file)
		if err = i.adaptors.Image.UploadImage(reader, fileHeader[0].Filename); err != nil {
			errs = append(errs, err)
		}

		file.Close()
	}

	return
}

func (i *ImageCommand) GetList(limit, offset int64, order, sortBy string) (items []*m.Image, total int64, err error) {

	items, total, err = i.adaptors.Image.List(limit, offset, order, sortBy)

	return
}