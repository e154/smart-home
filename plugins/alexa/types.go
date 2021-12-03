// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package alexa

import (
	m "github.com/e154/smart-home/models"
	"github.com/gin-gonic/gin"
)

const (
	// Name ...
	Name = "alexa"
	// TopicPluginAlexa ...
	TopicPluginAlexa = "plugin.alexa"
)

// ConfirmationStatus represents the status of either a dialog or slot confirmation.
type ConfirmationStatus string

const (
	// ConfConfirmed indicates the intent or slot has been confirmed by the end user.
	ConfConfirmed ConfirmationStatus = "CONFIRMED"

	// ConfDenied means the end user indicated the intent or slot should NOT proceed.
	ConfDenied ConfirmationStatus = "DENIED"

	// ConfNone means there has been not acceptance or denial of the intent or slot.
	ConfNone ConfirmationStatus = "NONE"
)

// Type will indicate type of dialog interaction to be sent to the user.
type DialogType string

const (
	// Delegate will indicate that the Server service should continue the dialog ineraction.
	Delegate DialogType = "Dialog.Delegate"

	// ElicitSlot will indicate to the Server service that the specific slot should be elicited from the user.
	ElicitSlot DialogType = "Dialog.ElicitSlot"

	// ConfirmSlot indicates to the Server service that the slot value should be confirmed by the user.
	ConfirmSlot DialogType = "Dialog.ConfirmSlot"

	// ConfirmIntent indicates to the Server service that the complete intent should be confimed by the user.
	ConfirmIntent DialogType = "Dialog.ConfirmIntent"
)

const (
	// Started indicates that the dialog interaction has just begun.
	Started string = "STARTED"

	// InProgress indicates that the dialog interation is continuing.
	InProgress string = "IN_PROGRESS"

	// Completed indicates that the dialog interaction has finished.
	// The intent and slot confirmation status should be checked.
	Completed string = "COMPLETED"
)

var (
	insecureSkipVerify = false
)

// ReqBody contains all data related to the type of request sent.
type ReqBody struct {
	Type        string `json:"type"`
	RequestID   string `json:"requestId"`
	Timestamp   string `json:"timestamp"`
	Intent      Intent `json:"intent,omitempty"`
	Reason      string `json:"reason,omitempty"`
	Locale      string `json:"locale,omitempty"`
	DialogState string `json:"dialogState,omitempty"`
}

// Intent represents the intent that is sent as part of an Request. This includes
// the name of the intent configured in the Server developers dashboard as well as any slots
// and the optional confirmation status if one is needed to complete an intent.
type Intent struct {
	Name               string             `json:"name"`
	Slots              map[string]Slot    `json:"slots"`
	ConfirmationStatus ConfirmationStatus `json:"confirmationStatus"`
}

// Slot represents variable values that can be sent that were specified by the end user
// when invoking the Server application.
type Slot struct {
	Name               string             `json:"name"`
	Value              string             `json:"value"`
	Resolutions        Resolution         `json:"resolutions"`
	ConfirmationStatus ConfirmationStatus `json:"confirmationStatus"`
}

// Resolution contains the results of entity resolutions when it relates to slots and how
// the values are resolved. The resolutions will be organized by authority, for custom slots
// the authority will be the custom slot type that was defined.
// Find more information here: https://developer.amazon.com/docs/custom-skills/define-synonyms-and-ids-for-slot-type-values-entity-resolution.html#intentrequest-changes
type Resolution struct {
	ResolutionsPerAuthority []ResolutionPerAuthority `json:"resolutionsPerAuthority"`
}

// ResolutionPerAuthority contains information about a single slot resolution from a single
// authority. The values silce will contain all possible matches for different slots.
// These resolutions are most interesting when working with synonyms.
type ResolutionPerAuthority struct {
	Authority string `json:"authority"`
	Status    struct {
		Code string `json:"code"`
	} `json:"status"`
	Values []map[string]struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"values"`
}

// Response represents the information that should be sent back to the Server service
// from the skillserver.
type Response struct {
	Version           string                 `json:"version"`
	SessionAttributes map[string]interface{} `json:"sessionAttributes,omitempty"`
	Response          RespBody               `json:"response"`
}

// RespBody contains the body of the response to be sent back to the Server service.
// This includes things like the text that should be spoken or any cards that should
// be shown in the Server companion app.
type RespBody struct {
	OutputSpeech     *RespPayload `json:"outputSpeech,omitempty"`
	Card             *RespPayload `json:"card,omitempty"`
	Reprompt         *Reprompt    `json:"reprompt,omitempty"` // Pointer so it's dropped if empty in JSON response.
	ShouldEndSession bool         `json:"shouldEndSession"`
	Directives       []*Directive `json:"directives,omitempty"`
}

// Reprompt contains speech that should be spoken back to the end user to retrieve
// additional information or to confirm an action.
type Reprompt struct {
	OutputSpeech RespPayload `json:"outputSpeech,omitempty"`
}

// EchoRespImage represents a single image with two variants that should be returned as part
// of a response. Small and Large image sizes can be provided.
type EchoRespImage struct {
	SmallImageURL string `json:"smallImageUrl,omitempty"`
	LargeImageURL string `json:"largeImageUrl,omitempty"`
}

// RespPayload contains the interesting parts of the Echo response including text to be spoken,
// card attributes, and images.
type RespPayload struct {
	Type    string        `json:"type,omitempty"`
	Title   string        `json:"title,omitempty"`
	Text    string        `json:"text,omitempty"`
	SSML    string        `json:"ssml,omitempty"`
	Content string        `json:"content,omitempty"`
	Image   EchoRespImage `json:"image,omitempty"`
}

// Directive includes information about intents and slots that should be confirmed or elicted from the user.
// The type value can be used to delegate the action to the Server service. In this case, a pre-configured prompt
// will be used from the developer console.
type Directive struct {
	Type            DialogType `json:"type"`
	UpdatedIntent   *Intent    `json:"updatedIntent,omitempty"`
	SlotToConfirm   string     `json:"slotToConfirm,omitempty"`
	SlotToElicit    string     `json:"slotToElicit,omitempty"`
	IntentToConfirm string     `json:"intentToConfirm,omitempty"`
}

// Session contains information about the ongoing session between the Server server and
// the skillserver. This session is stored as part of each request.
type Session struct {
	New         bool   `json:"new"`
	SessionID   string `json:"sessionId"`
	Application struct {
		ApplicationID string `json:"applicationId"`
	} `json:"application"`
	Attributes map[string]interface{} `json:"attributes"`
	User       struct {
		UserID      string `json:"userId"`
		AccessToken string `json:"accessToken,omitempty"`
	} `json:"user"`
}

// Context contains information about the context in which the request was sent.
// This could be information about the device from which the request was sent or about the invoked Server application.
type Context struct {
	System struct {
		Device struct {
			DeviceID string `json:"deviceId,omitempty"`
		} `json:"device,omitempty"`
		Application struct {
			ApplicationID string `json:"applicationId,omitempty"`
		} `json:"application,omitempty"`
		User struct {
			UserId string `json:"userId"`
		} `json:"user,omitempty"`
		ApiEndpoint    string `json:"apiEndpoint"`
		ApiAccessToken string `json:"apiAccessToken"`
	} `json:"System,omitempty"`
}

// EventAlexaAction ...
type EventAlexaAction struct {
	SkillId    int64
	IntentName string
	Payload    interface{}
}

// EventAlexaAddSkill ...
type EventAlexaAddSkill struct {
	Skill *m.AlexaSkill
}

// EventAlexaUpdateSkill ...
type EventAlexaUpdateSkill struct {
	Skill *m.AlexaSkill
}

// EventAlexaDeleteSkill ...
type EventAlexaDeleteSkill struct {
	Skill *m.AlexaSkill
}

// IServer ...
type IServer interface {
	Start()
	Stop()
	OnLaunchHandler(ctx *gin.Context, req *Request, resp *Response)
	OnIntentHandle(ctx *gin.Context, req *Request, resp *Response)
	OnSessionEndedHandler(ctx *gin.Context, req *Request, resp *Response)
	OnAudioPlayerHandler(ctx *gin.Context, req *Request, resp *Response)
	AddSkill(skill *m.AlexaSkill)
	UpdateSkill(skill *m.AlexaSkill)
	DeleteSkill(skill *m.AlexaSkill)
}

// AlexaPlugin ...
type AlexaPlugin interface {
	Server() IServer
}

const (
	// TriggerOptionSkillId ...
	TriggerOptionSkillId = "skillId"
)
