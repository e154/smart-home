// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package controllers

import (
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/stream"
	"os"
)

// ControllerImage ...
type ControllerImage struct {
	*ControllerCommon
}

// NewControllerImage ...
func NewControllerImage(common *ControllerCommon) *ControllerImage {
	return &ControllerImage{
		ControllerCommon: common,
	}
}

// Start ...
func (c *ControllerImage) Start() {
	c.stream.Subscribe("get_image_list", c.GetImageList)
	c.stream.Subscribe("get_filter_list", c.GetFilterList)
	c.stream.Subscribe("remove_image", c.RemoveImage)
}

// Stop ...
func (c *ControllerImage) Stop() {
	c.stream.UnSubscribe("get_image_list")
	c.stream.UnSubscribe("get_filter_list")
	c.stream.UnSubscribe("remove_image")
}

// Stream
func (c *ControllerImage) GetImageList(client stream.IStreamClient, message stream.Message) {
	//fmt.Println("get_image_list")

	filter, _ := message.Payload["filter"].(string)

	images, err := c.adaptors.Image.GetAllByDate(filter)
	if err != nil {
		client.Notify("error", err.Error())
		log.Error(err.Error())
		return
	}

	payload := map[string]interface{}{"images": images}
	response := message.Response(payload)
	client.Write(response.Pack())
}

// GetFilterList ...
func (c *ControllerImage) GetFilterList(client stream.IStreamClient, message stream.Message) {

	filterList, err := c.adaptors.Image.GetFilterList()
	if err != nil {
		client.Notify("error", err.Error())
		log.Error(err.Error())
		return
	}

	payload := map[string]interface{}{"filter_list": filterList}
	response := message.Response(payload)
	client.Write(response.Pack())
}

// RemoveImage ...
func (c *ControllerImage) RemoveImage(client stream.IStreamClient, message stream.Message) {

	fileId, ok := message.Payload["image_id"].(float64)
	if !ok {
		client.Notify("error", "image remove: bad image id request")
		log.Warn("image remove: bad image id request")
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

	client.Write(message.Success().Pack())
}
