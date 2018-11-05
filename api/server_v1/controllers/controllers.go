package controllers

import "github.com/e154/smart-home/adaptors"

type Controllers struct {
	Index *ControllerIndex
}

func NewControllers(adaptors *adaptors.Adaptors) *Controllers {
	common := NewControllerCommon(adaptors)
	return &Controllers{
		Index: NewControllerIndex(common),
	}
}
