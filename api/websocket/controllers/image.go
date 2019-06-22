package controllers

import "github.com/e154/smart-home/system/stream"

type ControllerImage struct {
	*ControllerCommon
}

func NewControllerImage(common *ControllerCommon,
	stream *stream.StreamService) *ControllerImage {
	image := &ControllerImage{
		ControllerCommon: common,
	}

	// register methods
	//stream.Subscribe("get_image_list", image.GetImageList)

	return image
}

func (c *ControllerImage) GetImageList(client *stream.Client, message stream.Message) {

	//server, err := c.GetServer(client)
	//if err != nil {
	//	c.Err(client, message, err)
	//	return
	//}
	//
	//log.Info("call register mobile")
	//
	//token, err := c.endpoint.RegisterImage(server)
	//if err != nil {
	//	c.Err(client, message, err)
	//	return
	//}
	//
	//client.Token = token
	//
	//payload := map[string]interface{}{
	//	"token": token,
	//}
	//response := message.Response(payload)
	//
	//client.Send <- response.Pack()

	return
}