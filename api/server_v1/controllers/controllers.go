package controllers

type Controllers struct {
	Index *ControllerIndex
}

func NewControllers() *Controllers {
	common := NewControllerCommon()
	return &Controllers{
		Index: NewControllerIndex(common),
	}
}
