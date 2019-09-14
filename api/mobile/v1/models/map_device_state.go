package models

import "encoding/json"

type MapDeviceState struct {
	Id            int64        `json:"id"`
	DeviceState   *DeviceState `json:"device_state"`
	DeviceStateId int64        `json:"device_state_id" valid:"Required"`
	MapDeviceId   int64        `json:"map_device_id" valid:"Required"`
	Image         *Image       `json:"image"`
	ImageId       int64        `json:"image_id" valid:"Required"`
	Style         string       `json:"style"`
}

func (n MapDeviceState) MarshalJSON() (b []byte, err error) {

	b, err = json.Marshal(map[string]interface{}{
		"id":              n.Id,
		"system_name":     n.DeviceState.SystemName,
		"description":     n.DeviceState.Description,
		"device_state_id": n.DeviceStateId,
		"device_id":       n.MapDeviceId,
		"image":           n.Image,
		"style":           n.Style,
	})

	return
}
