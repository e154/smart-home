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

package messagebird

import (
	"github.com/e154/smart-home/common"
	"github.com/messagebird/go-rest-api"
	"github.com/messagebird/go-rest-api/balance"
	"github.com/messagebird/go-rest-api/sms"
	"github.com/pkg/errors"
)

var (
	log = common.MustGetLogger("message bird")
)

// MBClient ...
type MBClient struct {
	cfg    *MBClientConfig
	client *messagebird.Client
}

// NewMBClient ...
func NewMBClient(cfg *MBClientConfig) (*MBClient, error) {

	if cfg.Name == "" || cfg.AccessKey == "" {
		return nil, errors.New("bad parameters")
	}

	client := &MBClient{
		cfg:    cfg,
		client: messagebird.New(cfg.AccessKey),
	}
	return client, nil
}

// SendSMS ...
func (c *MBClient) SendSMS(phone, body string) (string, error) {

	log.Infof("send sms %v, %v", phone, body)

	msgParams := &sms.Params{
		Type:       "sms",
		DataCoding: "unicode",
	}

	msg, err := sms.Create(c.client, c.cfg.Name, []string{phone}, body, msgParams)
	if err != nil {
		mbErr, ok := err.(messagebird.ErrorResponse)
		if !ok {
			return "", errors.New(err.Error())
		}

		//fmt.Println("Code:", mbErr.Errors[0].Code)
		//fmt.Println("Description:", mbErr.Errors[0].Description)
		//fmt.Println("Parameter:", mbErr.Errors[0].Parameter)

		return "", errors.New(mbErr.Errors[0].Description)
	}

	return msg.ID, nil
}

// GetStatus ...
func (c *MBClient) GetStatus(smsId string) (string, error) {

	msg, err := sms.Read(c.client, smsId)
	if err != nil {
		return "", errors.New(err.Error())
	}

	return msg.Recipients.Items[0].Status, nil
}

// Balance ...
func (c *MBClient) Balance() (*balance.Balance, error) {

	b, err := balance.Read(c.client)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return b, nil
}
