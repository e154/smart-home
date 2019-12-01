package responses

import "github.com/e154/smart-home/api/server/v1/models"

// swagger:response MqttClientList
type MqttClientList struct {
	// in:body
	Body struct {
		Items []*models.MqttClient `json:"items"`
		Meta  struct {
			Limit       int64 `json:"limit"`
			ObjectCount int64 `json:"object_count"`
			Offset      int64 `json:"offset"`
		} `json:"meta"`
	}
}

// swagger:response MqttSessionList
type MqttSessionList struct {
	// in:body
	Body struct {
		Items []*models.MqttSession `json:"items"`
		Meta  struct {
			Limit       int64 `json:"limit"`
			ObjectCount int64 `json:"object_count"`
			Offset      int64 `json:"offset"`
		} `json:"meta"`
	}
}

// swagger:response MqttSubscriptionList
type MqttSubscriptionList struct {
	// in:body
	Body struct {
		Items []*models.MqttSubscription `json:"items"`
		Meta  struct {
			Limit       int64 `json:"limit"`
			ObjectCount int64 `json:"object_count"`
			Offset      int64 `json:"offset"`
		} `json:"meta"`
	}
}
