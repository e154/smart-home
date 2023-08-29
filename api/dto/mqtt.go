package dto

import (
	"github.com/e154/smart-home/api/stub/api"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/mqtt/admin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Mqtt ...
type Mqtt struct{}

// NewMqttDto ...
func NewMqttDto() Mqtt {
	return Mqtt{}
}

// GetClientById ...
func (r Mqtt) GetClientById(from *admin.ClientInfo) (client *api.Client) {
	if from == nil {
		return
	}
	client = &api.Client{
		ClientId:             from.ClientID,
		Username:             from.Username,
		KeepAlive:            uint32(from.KeepAlive),
		Version:              from.Version,
		WillRetain:           from.WillRetain,
		WillQos:              uint32(from.WillQos),
		WillTopic:            from.WillTopic,
		WillPayload:          from.WillPayload,
		RemoteAddr:           from.RemoteAddr,
		LocalAddr:            from.LocalAddr,
		SubscriptionsCurrent: from.SubscriptionsCurrent,
		SubscriptionsTotal:   from.SubscriptionsTotal,
		PacketsReceivedBytes: from.PacketsReceivedBytes,
		PacketsReceivedNums:  from.PacketsReceivedNums,
		PacketsSendBytes:     from.PacketsSendBytes,
		PacketsSendNums:      from.PacketsSendNums,
		MessageDropped:       from.MessageDropped,
		InflightLen:          from.InflightLen,
		QueueLen:             from.QueueLen,
		ConnectedAt:          timestamppb.New(from.ConnectedAt),
	}
	if from.DisconnectedAt != nil {
		client.DisconnectedAt = timestamppb.New(*from.DisconnectedAt)
	}
	return
}

// ToListResult ...
func (r Mqtt) ToListResult(list []*admin.ClientInfo, total uint64, pagination common.PageParams) *api.GetClientListResult {

	items := make([]*api.Client, 0, len(list))

	for _, i := range list {
		items = append(items, r.GetClientById(i))
	}

	return &api.GetClientListResult{
		Items: items,
		Meta: &api.Meta{
			Limit: uint64(pagination.Limit),
			Page:  pagination.PageReq,
			Total: total,
			Sort:  pagination.SortReq,
		},
	}
}
