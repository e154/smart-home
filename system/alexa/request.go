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
	"errors"
	"time"
)

// Request represents all fields sent from the Alexa service to the skillserver.
// Convenience methods are provided to pull commonly used properties out of the request.
type Request struct {
	Version string  `json:"version"`
	Session Session `json:"session"`
	Context Context `json:"context"`
	Request ReqBody `json:"request"`
}

// VerifyTimestamp will parse the timestamp in the Request and verify that it is in the correct
// format and is not too old. True will be returned if the timestamp is valid; false otherwise.
func (r *Request) VerifyTimestamp() bool {
	reqTimestamp, _ := time.Parse("2006-01-02T15:04:05Z", r.Request.Timestamp)
	if time.Since(reqTimestamp) < time.Duration(150)*time.Second {
		return true
	}

	return false
}

// VerifyAppID check that the incoming application ID matches the application ID provided
// when running the server. This is a step required for skill certification.
func (r *Request) VerifyAppID(myAppID string) bool {
	if r.Session.Application.ApplicationID == myAppID ||
		r.Context.System.Application.ApplicationID == myAppID {
		return true
	}

	return false
}

// GetSessionID is a convenience method for getting the session ID out of an Request.
func (r *Request) GetSessionID() string {
	return r.Session.SessionID
}

// GetUserID is a convenience method for getting the user identifier out of an Request.
func (r *Request) GetUserID() string {
	return r.Session.User.UserID
}

// GetRequestType is a convenience method for getting the request type out of an Request.
func (r *Request) GetRequestType() string {
	return r.Request.Type
}

// GetIntentName is a convenience method for getting the intent name out of an Request.
func (r *Request) GetIntentName() string {
	if r.GetRequestType() == "IntentRequest" {
		return r.Request.Intent.Name
	}

	return r.GetRequestType()
}

// GetSlotValue is a convenience method for getting the value of the specified slot out of an Request
// as a string. An error is returned if a slot with that value is not found in the request.
func (r *Request) GetSlotValue(slotName string) (string, error) {
	slot, err := r.GetSlot(slotName)

	if err != nil {
		return "", err
	}

	return slot.Value, nil
}

// GetSlot will return an Slot from the Request with the given name.
func (r *Request) GetSlot(slotName string) (Slot, error) {
	if _, ok := r.Request.Intent.Slots[slotName]; ok {
		return r.Request.Intent.Slots[slotName], nil
	}

	return Slot{}, errors.New("slot name not found")
}

// AllSlots will return a map of all the slots in the Request mapped by their name.
func (r *Request) AllSlots() map[string]Slot {
	return r.Request.Intent.Slots
}

// Locale returns the locale specified in the request.
func (r *Request) Locale() string {
	return r.Request.Locale
}
