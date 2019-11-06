package endpoint

import (
	m "github.com/e154/smart-home/models"
)

type MessageDeliveryEndpoint struct {
	*CommonEndpoint
}

func NewMessageDeliveryEndpoint(common *CommonEndpoint) *MessageDeliveryEndpoint {
	return &MessageDeliveryEndpoint{
		CommonEndpoint: common,
	}
}

func (n *MessageDeliveryEndpoint) GetList(limit, offset int64, order, sortBy string) (result []*m.MessageDelivery, total int64, err error) {
	result, total, err = n.adaptors.MessageDelivery.List(limit, offset, order, sortBy)
	return
}

func (n *MessageDeliveryEndpoint) Delete(id int64) (err error) {
	err = n.adaptors.MessageDelivery.Delete(id)
	return
}
