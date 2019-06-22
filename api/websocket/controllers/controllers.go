package controllers

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/endpoint"
	"github.com/e154/smart-home/system/stream"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("stream.controllers")
)

type Controllers struct {
	Image *ControllerImage
}

func NewControllers(adaptors *adaptors.Adaptors,
	stream *stream.StreamService,
	endpoint *endpoint.Endpoint) *Controllers {
	common := NewControllerCommon(adaptors, stream, endpoint)
	return &Controllers{
		Image: NewControllerImage(common, stream),
	}
}
