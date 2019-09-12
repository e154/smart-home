package models

import "encoding/json"

type MapDeviceAction struct {
	Id             int64         `json:"id"`
	DeviceAction   *DeviceAction `json:"device_action"`
	DeviceActionId int64         `json:"device_action_id" valid:"Required"`
	MapDeviceId    int64         `json:"map_device_id" valid:"Required"`
	Image          *Image        `json:"image"`
	ImageId        int64         `json:"image_id" valid:"Required"`
	Type           string        `json:"type"`
}

func (n MapDeviceAction) MarshalJSON() (b []byte, err error) {

	b, err = json.Marshal(map[string]interface{}{
		"id":               n.Id,
		"name":             n.DeviceAction.Name,
		"description":      n.DeviceAction.Description,
		"device_action_id": n.DeviceActionId,
		"image":            n.Image,
		"device_id":        n.MapDeviceId,
	})

	return
}
