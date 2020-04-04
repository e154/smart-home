// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/alexa/dialog"
	"github.com/gin-gonic/gin"
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

// Application represents a single Alexa application server. This application type needs to include
// the application ID from the Alexa developer portal that will be making requests to the server. This AppId needs
// to be verified to ensure the requests are coming from the correct app. Handlers can also be provied for
// different types of requests sent by the Alexa Skills Kit such as OnLaunch or OnIntent.
type Application interface {
	GetAppID() string
	OnLaunch(*gin.Context, *Request, *Response)
	OnIntent(*gin.Context, *Request, *Response)
	OnSessionEnded(*gin.Context, *Request, *Response)
	OnAudioPlayerState(*gin.Context, *Request, *Response)
}

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
// the name of the intent configured in the Alexa developers dashboard as well as any slots
// and the optional confirmation status if one is needed to complete an intent.
type Intent struct {
	Name               string             `json:"name"`
	Slots              map[string]Slot    `json:"slots"`
	ConfirmationStatus ConfirmationStatus `json:"confirmationStatus"`
}

// Slot represents variable values that can be sent that were specified by the end user
// when invoking the Alexa application.
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

// Response Types

// Response represents the information that should be sent back to the Alexa service
// from the skillserver.
type Response struct {
	Version           string                 `json:"version"`
	SessionAttributes map[string]interface{} `json:"sessionAttributes,omitempty"`
	Response          RespBody               `json:"response"`
}

// RespBody contains the body of the response to be sent back to the Alexa service.
// This includes things like the text that should be spoken or any cards that should
// be shown in the Alexa companion app.
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
// The type value can be used to delegate the action to the Alexa service. In this case, a pre-configured prompt
// will be used from the developer console.
type Directive struct {
	Type            dialog.Type `json:"type"`
	UpdatedIntent   *Intent     `json:"updatedIntent,omitempty"`
	SlotToConfirm   string      `json:"slotToConfirm,omitempty"`
	SlotToElicit    string      `json:"slotToElicit,omitempty"`
	IntentToConfirm string      `json:"intentToConfirm,omitempty"`
}

// Session contains information about the ongoing session between the Alexa server and
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
// This could be information about the device from which the request was sent or about the invoked Alexa application.
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

type Options func(a *Alexa)

func WithServerOption(addressPort string) Options {
	return func(a *Alexa) {
		a.addressPort = common.String(addressPort)
	}
}
