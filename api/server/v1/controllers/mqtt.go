package controllers

import (
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/e154/smart-home/common"
	"github.com/gin-gonic/gin"
)

type ControllerMqtt struct {
	*ControllerCommon
}

func NewControllerMqtt(common *ControllerCommon) *ControllerMqtt {
	return &ControllerMqtt{ControllerCommon: common}
}

// swagger:operation GET /mqtt/clients mqttClientList
// ---
// summary: get client list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - mqtt
// parameters:
// - default: 10
//   description: limit
//   in: query
//   name: limit
//   required: true
//   type: integer
// - default: 0
//   description: offset
//   in: query
//   name: offset
//   required: true
//   type: integer
// responses:
//   "200":
//	   $ref: '#/responses/MqttClientList'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerMqtt) GetClients(ctx *gin.Context) {

	_, _, _, limit, offset := c.list(ctx)
	items, total, err := c.endpoint.Mqtt.GetClients(limit, offset)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.MqttClient, 0)
	_ = common.Copy(&result, &items, common.JsonEngine)

	resp := NewSuccess()
	resp.Page(limit, offset, int64(total), result).Send(ctx)
}

// swagger:operation GET /mqtt/client/{id} mqttClientGetById
// ---
// parameters:
// - description: client ID
//   in: path
//   name: id
//   required: true
//   type: string
// summary: get client by ID
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - mqtt
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/MqttClient'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "404":
//	   $ref: '#/responses/Error'
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerMqtt) GetClientById(ctx *gin.Context) {

	id := ctx.Param("id")
	mqttClient, err := c.endpoint.Mqtt.GetClient(id)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := &models.MqttClient{}
	common.Copy(&result, &mqttClient, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /mqtt/sessions mqttSessionList
// ---
// summary: get session list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - mqtt
// parameters:
// - default: 10
//   description: limit
//   in: query
//   name: limit
//   required: true
//   type: integer
// - default: 0
//   description: offset
//   in: query
//   name: offset
//   required: true
//   type: integer
// responses:
//   "200":
//	   $ref: '#/responses/MqttSessionList'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerMqtt) GetSessions(ctx *gin.Context) {

	_, _, _, limit, offset := c.list(ctx)
	items, total, err := c.endpoint.Mqtt.GetSessions(limit, offset)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.MqttSession, 0)
	_ = common.Copy(&result, &items, common.JsonEngine)

	resp := NewSuccess()
	resp.Page(limit, offset, int64(total), result).Send(ctx)
}

// swagger:operation GET /mqtt/client/{id}/session mqttClientGetById
// ---
// parameters:
// - description: client ID
//   in: path
//   name: id
//   required: true
//   type: string
// summary: get session by client ID
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - mqtt
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/MqttSession'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "404":
//	   $ref: '#/responses/Error'
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerMqtt) GetSession(ctx *gin.Context) {

	id := ctx.Param("id")
	mqttClient, err := c.endpoint.Mqtt.GetSession(id)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := &models.MqttSession{}
	common.Copy(&result, &mqttClient, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /mqtt/client/{id}/subscriptions mqttSubscriptionList
// ---
// summary: get subscription list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - mqtt
// parameters:
// - description: client ID
//   in: path
//   name: id
//   required: true
//   type: string
// - default: 10
//   description: limit
//   in: query
//   name: limit
//   required: true
//   type: integer
// - default: 0
//   description: offset
//   in: query
//   name: offset
//   required: true
//   type: integer
// responses:
//   "200":
//	   $ref: '#/responses/MqttSubscriptionList'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerMqtt) GetSubscriptions(ctx *gin.Context) {

	id := ctx.Param("id")
	_, _, _, limit, offset := c.list(ctx)
	items, total, err := c.endpoint.Mqtt.GetSubscriptions(id, limit, offset)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.MqttSubscription, 0)
	_ = common.Copy(&result, &items, common.JsonEngine)

	resp := NewSuccess()
	resp.Page(limit, offset, int64(total), result).Send(ctx)
}

// swagger:operation DELETE /mqtt/client/{id}/topic mqttUnsubscribeTopic
// ---
// parameters:
// - description: client ID
//   in: path
//   name: id
//   required: true
//   type: string
// - description: topic
//   in: query
//   name: topic
//   required: true
//   type: string
// summary: delete mqtt by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - mqtt
// responses:
//   "200":
//	   $ref: '#/responses/Success'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "404":
//	   $ref: '#/responses/Error'
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerMqtt) Unsubscribe(ctx *gin.Context) {

	id := ctx.Param("id")
	topic := ctx.Query("topic")
	err := c.endpoint.Mqtt.Unsubscribe(id, topic)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	NewSuccess().Send(ctx)
	return
}

// swagger:operation POST /mqtt/publish mqttPublish
// ---
// parameters:
// - description: publish params
//   in: body
//   name: mqtt
//   required: true
//   schema:
//     $ref: '#/definitions/NewMqttPublish'
//     type: object
// summary: publish
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - mqtt
// responses:
//   "200":
//	   $ref: '#/responses/Success'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerMqtt) Publish(ctx *gin.Context) {

	params := &models.NewMqttPublish{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	err := c.endpoint.Mqtt.Publish(params.Topic, params.Qos, params.Payload, params.Retain)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	NewSuccess().Send(ctx)
}

// swagger:operation DELETE /mqtt/client/{id} mqttCloseClient
// ---
// parameters:
// - description: client ID
//   in: path
//   name: id
//   required: true
//   type: string
// summary: close client
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - mqtt
// responses:
//   "200":
//	   $ref: '#/responses/Success'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "404":
//	   $ref: '#/responses/Error'
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerMqtt) CloseClient(ctx *gin.Context) {

	id := ctx.Param("id")
	err := c.endpoint.Mqtt.CloseClient(id)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	NewSuccess().Send(ctx)
	return
}