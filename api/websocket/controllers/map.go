package controllers

import (
	"github.com/e154/smart-home/system/stream"
)

type ControllerMap struct {
	*ControllerCommon
}

func NewControllerMap(common *ControllerCommon) *ControllerMap {
	return &ControllerMap{
		ControllerCommon: common,
	}
}

func (c *ControllerMap) Start() {

}

func (c *ControllerMap) Stop() {

}

func (n *ControllerMap) NodesStatus(client *stream.Client, value interface{}) {
	//v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
	//if !ok {
	//	return
	//}
	//
	//n.Update()
	//
	//msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "nodes": n})
	//client.Send <- msg
}

