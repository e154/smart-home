package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"../models"
	"../stream"
	"reflect"
	"fmt"
	"../log"
	"os"
	"../../lib/common"
)

// ImageController operations for Image
type ImageController struct {
	CommonController
}

// URLMapping ...
func (c *ImageController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Post", c.Upload)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Image
// @Param	body		body 	models.Image	true		"body for Image content"
// @Success 201 {object} models.Image
// @Failure 403 body is empty
// @router / [post]
func (c *ImageController) Post() {
	var image models.Image
	json.Unmarshal(c.Ctx.Input.RequestBody, &image)

	// validation
	valid := validation.Validation{}
	b, err := valid.Valid(&image)
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	if !b {
		var msg string
		for _, err := range valid.Errors {
			msg += fmt.Sprintf("%s: %s\r", err.Key, err.Message)
		}
		c.ErrHan(403, msg)
		return
	}
	//....

	_, err = models.AddImage(&image)
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{"image": image}
	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get Image by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Image
// @Failure 403 :id is empty
// @router /:id [get]
func (c *ImageController) GetOne() {
	id, _ := c.GetInt(":id")
	image, err := models.GetImageById(int64(id))
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{"image": image}
	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get Image
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Image
// @Failure 403
// @router / [get]
func (c *ImageController) GetAll() {
	images, meta, err := models.GetAllImage(c.pagination())
	if err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.Data["json"] = &map[string]interface{}{"images": images, "meta": meta}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Image
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Image	true		"body for Image content"
// @Success 200 {object} models.Image
// @Failure 403 :id is not int
// @router /:id [put]
func (c *ImageController) Put() {
	id, _ := c.GetInt(":id")
	var image models.Image
	json.Unmarshal(c.Ctx.Input.RequestBody, &image)
	image.Id = int64(id)
	if err := models.UpdateImageById(&image); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Image
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *ImageController) Delete() {
	id, _ := c.GetInt(":id")
	if err := models.DeleteImage(int64(id)); err != nil {
		c.ErrHan(403, err.Error())
		return
	}

	c.ServeJSON()
}

func (c *ImageController) Upload() {

	// этот кусок не работал
	//qwe, err := c.GetFiles("files")
	files := c.Ctx.Request.MultipartForm.File
	if len(files) == 0 {
		c.ErrHan(403, "http: no such file")
		return
	}

	images, errs := models.UploadImages(files)

	c.Data["json"] = &map[string]interface{}{"images": images, "errors": errs}
	c.ServeJSON()
}

// Stream
func getImageList(client *stream.Client, value interface{}) {
	v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
	if !ok {
		return
	}

	var filter string
	if filter, ok = v["filter"].(string); ok {
	}

	ml, err := models.GetAllImagesByDate(filter)
	if err != nil {
		client.Notify("error", err.Error())
		log.Error(err.Error())
		return
	}

	images := []*models.Image{}
	for _, image := range ml {
		newImage := &models.Image{}
		*newImage = image
		newImage.Url = common.GetLinkPath(newImage.Image)
		images = append(images, newImage)
	}


	msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "images": images})
	client.Send(string(msg))
}

func getFilterList(client *stream.Client, value interface{}) {
	v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
	if !ok {
		return
	}

	filter_list, err := models.GetImageFilterList()
	if err != nil {
		client.Notify("error", err.Error())
		log.Error(err.Error())
		return
	}

	msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "filter_list": filter_list})
	client.Send(string(msg))
}

func removeImage(client *stream.Client, value interface{}) {
	v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
	if !ok {
		return
	}

	var file_id float64
	if file_id, ok = v["image_id"].(float64); !ok {
		client.Notify("error", "image remove: bad image id request")
		log.Warn("image remove: bad image id request")
		return
	}

	image, err := models.GetImageById(int64(file_id))
	if err != nil {
		client.Notify("error", err.Error())
		log.Error(err.Error())
		return
	}

	if err = models.DeleteImage(image.Id); err != nil {
		client.Notify("error", err.Error())
		return
	}

	file_path := common.GetFullPath(image.Image)
	if err = os.Remove(file_path); err != nil {
		log.Warnf("remove image file:",err.Error())
	}

	msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "status": "ok"})
	client.Send(string(msg))
}

func init() {
	hub := stream.GetHub()
	hub.Subscribe("get_image_list", getImageList)
	hub.Subscribe("get_filter_list", getFilterList)
	hub.Subscribe("remove_image", removeImage)
}