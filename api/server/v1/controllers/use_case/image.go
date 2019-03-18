package use_case

import (
	"errors"
	"mime/multipart"
	"github.com/jinzhu/copier"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/validation"
	"github.com/e154/smart-home/api/server/v1/models"
	m "github.com/e154/smart-home/models"
	"bufio"
)

func AddImage(imageParams models.NewImage, adaptors *adaptors.Adaptors) (ok bool, image *m.Image, errs []*validation.Error, err error) {

	image = &m.Image{}

	if err = copier.Copy(&image, &imageParams); err != nil {
		return
	}

	// validation
	ok, errs = image.Valid()
	if len(errs) > 0 || !ok {
		return
	}

	var id int64
	if id, err = adaptors.Image.Add(image); err != nil {
		return
	}

	image.Id = id

	return
}

func GetImageById(imageId int64, adaptors *adaptors.Adaptors) (image *m.Image, err error) {

	image, err = adaptors.Image.GetById(imageId)

	return
}

func UpdateImage(imageParams *models.UpdateImage, adaptors *adaptors.Adaptors) (result *models.Image, errs []*validation.Error, err error) {

	var image *m.Image
	if image, err = adaptors.Image.GetById(imageParams.Id); err != nil {
		return
	}

	if err = copier.Copy(&image, &imageParams); err != nil {
		return
	}

	// validation
	_, errs = image.Valid()
	if len(errs) > 0 {
		return
	}

	if err = adaptors.Image.Update(image); err != nil {
		return
	}

	if image, err = adaptors.Image.GetById(imageParams.Id); err != nil {
		return
	}

	result = &models.Image{}
	err = copier.Copy(&result, &image)

	return
}

func DeleteImageById(imageId int64, adaptors *adaptors.Adaptors) (err error) {

	if imageId == 0 {
		err = errors.New("image id is null")
		return
	}

	var image *m.Image
	if image, err = adaptors.Image.GetById(imageId); err != nil {
		return
	}

	err = adaptors.Image.Delete(image.Id)

	return
}

func UploadImages(files map[string][]*multipart.FileHeader, adaptors *adaptors.Adaptors) (fileList []*m.Image, errs []error) {

	fileList = make([]*m.Image, 0)
	errs = make([]error, 0)

	for _, fileHeader := range files {

		file, err := fileHeader[0].Open()
		if err != nil {
			errs = append(errs, err)
			continue
		}

		reader := bufio.NewReader(file)
		if err = adaptors.Image.UploadImage(reader, fileHeader[0].Filename); err != nil {
			errs = append(errs, err)
		}

		file.Close()
	}

	return
}

func GetImageList(limit, offset int64, order, sortBy string, adaptors *adaptors.Adaptors) (items []*m.Image, total int64, err error) {

	items, total, err = adaptors.Image.List(limit, offset, order, sortBy)

	return
}