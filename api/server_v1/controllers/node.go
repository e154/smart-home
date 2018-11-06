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

	node, err := GetNodeById(params.ID, c.adaptors)
	if err != nil {
		return err
	}

	resp := NewSuccess()
	resp.Item("node", node)
	return resp
}

func (c ControllerNode) UpdateNode(params node.UpdateNodeParams, principal interface{}) middleware.Responder {

	n := m.NewNode()
	n.Id = params.ID
	n.Port = int(params.PutNodeParams.Port)
	n.Status = params.PutNodeParams.Status
	n.Name = params.PutNodeParams.Name
	n.Ip = params.PutNodeParams.IP
	n.Description = params.PutNodeParams.Description

	if err := UpdateNode(n, c.adaptors, c.core); err != nil {
		return err
	}

	resp := NewSuccess()
	return resp
}

func (c ControllerNode) GetNodeList(params node.GetNodeListParams, principal interface{}) middleware.Responder {

	var limit int64 = 15
	var offset int64 = 0
	order := "DESC"
	sortBy := "id"

	if params.Limit != nil {
		limit = *params.Limit
	}
	if params.Limit != nil {
		offset = *params.Offset
	}
	if params.Order != nil {
		order = *params.Order
	}
	if params.Sortby != nil {
		order = *params.Sortby
	}

	items, total, err := GetNodeList(limit, offset, order, sortBy, c.adaptors)
	if err != nil {
		return err
	}

	resp := NewSuccess()
	resp.Page(limit, offset, total, items)
	return resp
}

func (c ControllerNode) DeleteNodeById(params node.DeleteNodeByIDParams, principal interface{}) middleware.Responder {

	if err := DeleteNodeById(params.ID, c.adaptors, c.core); err != nil {
		return err
	}

	resp := NewSuccess()
	return resp
}
