package controllers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/e154/smart-home/api/server_v1/stub/restapi/operations/node"
	. "github.com/e154/smart-home/api/server_v1/controllers/use_case"
	m "github.com/e154/smart-home/models"
)

type ControllerNode struct {
	*ControllerCommon
}

func NewControllerNode(common *ControllerCommon) *ControllerNode {
	return &ControllerNode{ControllerCommon: common}
}

func (c ControllerNode) AddNode(params node.AddNodeParams, principal interface{}) middleware.Responder {

	n := m.NewNode()
	n.Port = int(params.PostNodeParams.Port)
	n.Status = params.PostNodeParams.Status
	n.Name = params.PostNodeParams.Name
	n.Ip = params.PostNodeParams.IP
	n.Description = params.PostNodeParams.Description


	id, err := AddNode(n, c.adaptors, c.core)
	if err != nil {
		return err
	}

	resp := NewSuccess()
	resp.Item("id", id)
	return resp
}

func (c ControllerNode) GetNodeById(params node.GetNodeByIDParams, principal interface{}) middleware.Responder {

	resp := NewSuccess()
	return resp
}

func (c ControllerNode) GetNodeList(params node.GetNodeListParams, principal interface{}) middleware.Responder {

	resp := NewSuccess()
	return resp
}

func (c ControllerNode) DeleteNodeById(params node.DeleteNodeByIDParams, principal interface{}) middleware.Responder {

	resp := NewSuccess()
	return resp
}
