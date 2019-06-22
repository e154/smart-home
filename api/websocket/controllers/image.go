package controllers

import (
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/stream"
	"os"
)

type ControllerImage struct {
	*ControllerCommon
}

func NewControllerImage(common *ControllerCommon, ) *ControllerImage {
	return &ControllerImage{
		ControllerCommon: common,
	}
}

func (c *ControllerImage) Start() {
	c.stream.Subscribe("get_image_list", c.GetImageList)
	c.stream.Subscribe("get_filter_list", c.GetFilterList)
	c.stream.Subscribe("remove_image", c.RemoveImage)
}

func (c *ControllerImage) Stop() {
	c.stream.UnSubscribe("get_image_list")
	c.stream.UnSubscribe("get_filter_list")
	c.stream.UnSubscribe("remove_image")
}

// Stream
func (c *ControllerImage) GetImageList(client *stream.Client, message stream.Message) {
	//fmt.Println("get_image_list")

	filter, _ := message.Payload["filter"].(string)

	images, err := c.adaptors.Image.GetAllByDate(filter)
	if err != nil {
		client.Notify("error", err.Error())
		log.Error(err.Error())
		return
	}

	payload := map[string]interface{}{"images": images,}
	response := message.Response(payload)
	client.Send <- response.Pack()
}

func (c *ControllerImage) GetFilterList(client *stream.Client, message stream.Message) {

	filterList, err := c.adaptors.Image.GetFilterList()
	if err != nil {
		client.Notify("error", err.Error())
		log.Error(err.Error())
		return
	}

	payload := map[string]interface{}{"filter_list": filterList,}
	response := message.Response(payload)
	client.Send <- response.Pack()
}

func (c *ControllerImage) RemoveImage(client *stream.Client, message stream.Message) {

	fileId, ok := message.Payload["image_id"].(float64)
	if !ok {
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
	_ = os.Remove(filePath)

	client.Send <- message.Success().Pack()
}
