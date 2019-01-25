package controllers

import (
	"github.com/op/go-logging"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/stream"
)

var (
	log = logging.MustGetLogger("stream.controllers")
)

type StreamControllers struct {
	Image *ControllerImage
}

func NewStreamControllers(adaptors *adaptors.Adaptors,
	stream *stream.StreamService) *StreamControllers {

	ctrls := &StreamControllers{
		Image: NewControllerImage(adaptors),
	}

	// images
	stream.Subscribe("get_image_list", ctrls.Image.GetImageList)
	stream.Subscribe("get_filter_list", ctrls.Image.GetFilterList)
	stream.Subscribe("remove_image", ctrls.Image.RemoveImage)

	return ctrls
}
