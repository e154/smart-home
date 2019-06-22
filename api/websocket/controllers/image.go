package controllers

import (
	"encoding/json"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/stream"
	"os"
	"reflect"
)

type ControllerImage struct {
	*ControllerCommon
}

func NewControllerImage(common *ControllerCommon,) *ControllerImage {
	image := &ControllerImage{
		ControllerCommon: common,
	}

	// register methods
	common.stream.Subscribe("get_image_list", image.GetImageList)
	common.stream.Subscribe("get_filter_list", image.GetFilterList)
	common.stream.Subscribe("remove_image", image.RemoveImage)

	return image
}

// Stream
func (c *ControllerImage) GetImageList(client *stream.Client, value interface{}) {
	//fmt.Println("get_image_list")

	v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
	if !ok {
		return
	}

	var filter string
	if filter, ok = v["filter"].(string); ok {
	}

	images, err := c.adaptors.Image.GetAllByDate(filter)
	if err != nil {
		client.Notify("error", err.Error())
		log.Error(err.Error())
		return
	}

	msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "images": images})
	client.Send <- msg
}

func (c *ControllerImage) GetFilterList(client *stream.Client, value interface{}) {
	v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
	if !ok {
		return
	}

	filterList, err := c.adaptors.Image.GetFilterList()
	if err != nil {
		client.Notify("error", err.Error())
		log.Error(err.Error())
		return
	}

	msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "filter_list": filterList})
	client.Send <- msg
}

func (c *ControllerImage) RemoveImage(client *stream.Client, value interface{}) {
	v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
	if !ok {
		return
	}

	var fileId float64
	if fileId, ok = v["image_id"].(float64); !ok {
		client.Notify("error", "image remove: bad image id request")
		log.Warning("image remove: bad image id request")
		return
	}

	image, err := c.adaptors.Image.GetById(int64(fileId))
	if err != nil {
		client.Notify("error", err.Error())
		log.Error(err.Error())
		return
	}

	if err = c.adaptors.Image.Delete(image.Id); err != nil {
		client.Notify("error", err.Error())
		return
	}

	filePath := common.GetFullPath(image.Image)
	os.Remove(filePath)

	msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "status": "ok"})
	client.Send <- msg
}
