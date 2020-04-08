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
	"encoding/json"
)

// NewEchoResponse will construct a new response instance with the required metadata and an empty speech string.
// By default the response will indicate that the session should be ended. Use the `EndSession(bool)` method if the
// session should be left open.
func NewResponse() *Response {
	er := &Response{
		Version: "1.0",
		Response: RespBody{
			ShouldEndSession: true,
		},
		SessionAttributes: make(map[string]interface{}),
	}

	return er
}

// OutputSpeech will replace any existing text that should be spoken with this new value. If the output
// needs to be constructed in steps or special speech tags need to be used, see the `SSMLTextBuilder`.
func (r *Response) OutputSpeech(text string) *Response {
	r.Response.OutputSpeech = &RespPayload{
		Type: "PlainText",
		Text: text,
	}

	return r
}

// Card will add a card to the Alexa app's response with the provided title and content strings.
func (r *Response) Card(title string, content string) *Response {
	return r.SimpleCard(title, content)
}

// OutputSpeechSSML will add the text string provided and indicate the speech type is SSML in the response.
// This should only be used if the text to speech string includes special SSML tags.
func (r *Response) OutputSpeechSSML(text string) *Response {
	r.Response.OutputSpeech = &RespPayload{
		Type: "SSML",
		SSML: text,
	}

	return r
}

// SimpleCard will indicate that a card should be included in the Alexa companion app as part of the response.
// The card will be shown with the provided title and content.
func (r *Response) SimpleCard(title string, content string) *Response {
	r.Response.Card = &RespPayload{
		Type:    "Simple",
		Title:   title,
		Content: content,
	}

	return r
}

// StandardCard will indicate that a card should be shown in the Alexa companion app as part of the response.
// The card shown will include the provided title and content as well as images loaded from the locations provided
// as remote locations.
func (r *Response) StandardCard(title string, content string, smallImg string, largeImg string) *Response {
	r.Response.Card = &RespPayload{
		Type:    "Standard",
		Title:   title,
		Content: content,
	}

	if smallImg != "" {
		r.Response.Card.Image.SmallImageURL = smallImg
	}

	if largeImg != "" {
		r.Response.Card.Image.LargeImageURL = largeImg
	}

	return r
}

// LinkAccountCard is used to indicate that account linking still needs to be completed to continue
// using the Alexa skill. This will force an account linking card to be shown in the user's companion app.
func (r *Response) LinkAccountCard() *Response {
	r.Response.Card = &RespPayload{
		Type: "LinkAccount",
	}

	return r
}

// Reprompt will send a prompt back to the user, this could be used to request additional information from the user.
func (r *Response) Reprompt(text string) *Response {
	r.Response.Reprompt = &Reprompt{
		OutputSpeech: RespPayload{
			Type: "PlainText",
			Text: text,
		},
	}

	return r
}

// RepromptSSML is similar to the `Reprompt` method but should be used when the prompt
// to the user should include special speech tags.
func (r *Response) RepromptSSML(text string) *Response {
	r.Response.Reprompt = &Reprompt{
		OutputSpeech: RespPayload{
			Type: "SSML",
			Text: text,
		},
	}

	return r
}

// EndSession is a convenience method for setting the flag in the response that will
// indicate if the session between the end user's device and the skillserver should be closed.
func (r *Response) EndSession(flag bool) *Response {
	r.Response.ShouldEndSession = flag

	return r
}

// RespondToIntent is used to Delegate/Elicit/Confirm a dialog or an entire intent with
// user of alexa. The func takes in name of the dialog, updated intent/intent to confirm
// if any and optional slot value. It prepares a Echo Response to be returned.
// Multiple directives can be returned by calling the method in chain
// (eg. RespondToIntent(...).RespondToIntent(...), each RespondToIntent call appends the
// data to Directives array and will return the same at the end.
func (r *Response) RespondToIntent(name DialogType, intent *Intent, slot *Slot) *Response {
	directive := Directive{Type: name}
	if intent != nil && name == ConfirmIntent {
		directive.IntentToConfirm = intent.Name
	} else {
		directive.UpdatedIntent = intent
	}

	if slot != nil {
		if name == ElicitSlot {
			directive.SlotToElicit = slot.Name
		} else if name == ConfirmSlot {
			directive.SlotToConfirm = slot.Name
		}
	}
	r.Response.Directives = append(r.Response.Directives, &directive)
	return r
}

// String ...
func (r *Response) String() ([]byte, error) {
	jsonStr, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return jsonStr, nil
}
