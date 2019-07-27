package controllers

import (
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/stream"
)

type ControllerWorker struct {
	*ControllerCommon
}

func NewControllerWorker(common *ControllerCommon) *ControllerWorker {
	return &ControllerWorker{
		ControllerCommon: common,
	}
}

func (c *ControllerWorker) Start() {
	c.stream.Subscribe("do.worker", c.DoWorker)
}

func (c *ControllerWorker) Stop() {
	c.stream.UnSubscribe("do.worker")
}

// Stream
func (c *ControllerWorker) DoWorker(client *stream.Client, message stream.Message) {

	v := message.Payload
	var ok bool

	var workerId float64
	var err error

	if workerId, ok = v["worker_id"].(float64); !ok {
		log.Warning("bad id param")
		return
	}

	var worker *m.Worker
	if worker, err = c.adaptors.Worker.GetById(int64(workerId)); err != nil {
		client.Notify("error", err.Error())
		return
	}

	if err = c.core.DoWorker(worker); err != nil {
		client.Notify("error", err.Error())
		return
	}

	client.Send <- message.Success().Pack()
}
