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

package bitmine

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/e154/smart-home/common"
)

var (
	log = common.MustGetLogger("plugins.cgminer.bitmine")
)

// Bitmine ...
type Bitmine struct {
	device     DeviceType
	transport  ITransport
	user, pass string
}

// NewBitmine ...
func NewBitmine(transport ITransport, device, user, pass string) (bitmine *Bitmine, err error) {

	bitmine = &Bitmine{
		transport: transport,
		user:      url.QueryEscape(user),
		pass:      url.QueryEscape(pass),
	}

	switch device {
	case DeviceS9.String():
		bitmine.device = DeviceS9
	case DeviceS7.String():
		bitmine.device = DeviceS7
	case DeviceL3.String(), DeviceL3Plus.String():
		bitmine.device = DeviceL3Plus
	case DeviceD3.String():
		bitmine.device = DeviceD3
	case DeviceT9.String():
		bitmine.device = DeviceT9
	default:
		bitmine = nil
		err = fmt.Errorf("unknown device %s", device)
	}

	return
}

func (b *Bitmine) checkStatus(statuses []Status) error {
	for _, status := range statuses {
		switch status.Status {
		case "E":
			return fmt.Errorf("API returned error: Code: %d, Msg: '%s', Description: '%s'", status.Code, status.Msg, status.Description)
		case "F":
			return fmt.Errorf("API returned FATAL error: Code: %d, Msg: '%s', Description: '%s'", status.Code, status.Msg, status.Description)
		case "S":
		case "W":
		case "I":

		}
	}
	return nil
}

// Stats ...
func (b *Bitmine) Stats() (data []byte, err error) {

	var commandResponse []byte
	if commandResponse, err = b.transport.RunCommand("stats", ""); err != nil {
		return
	}

	var resp StatsResponse
	//TODO uncomment in go2
	//switch b.device {
	//case DeviceS9:
	//	resp = StatsResponse[StatsS9]{}
	//case DeviceS7:
	//	resp = StatsResponse[StatsS7]{}
	//case DeviceL3, DeviceL3Plus:
	//	resp = StatsResponse[StatsL3]{}
	//case DeviceD3:
	//	resp = StatsResponse[StatsD3]{}
	//case DeviceT9:
	//	resp = StatsResponse[StatsT9]{}
	//}

	// fix incorrect json response from miner "}{"
	fixResponse := bytes.Replace(commandResponse, []byte("}{"), []byte(","), 1)
	if err = json.Unmarshal(fixResponse, &resp); err != nil {
		return
	}
	if err = b.checkStatus(resp.Status); err != nil {
		return
	}
	if len(resp.Stats) < 1 {
		err = errors.New("no stats in JSON response")
		return
	}
	if len(resp.Stats) > 1 {
		err = errors.New("too many stats in JSON response")
		return
	}
	data, err = json.Marshal(resp.Stats[0])
	return
}

// Devs ...
func (b *Bitmine) Devs() (data []byte, err error) {
	var commandResponse []byte
	if commandResponse, err = b.transport.RunCommand("devs", ""); err != nil {
		return
	}
	var resp DevsResponse
	if err = json.Unmarshal(commandResponse, &resp); err != nil {
		return
	}
	if err = b.checkStatus(resp.Status); err != nil {
		return
	}
	data, err = json.Marshal(resp.Devs)
	return
}

// Summary ...
func (b *Bitmine) Summary() (data []byte, err error) {
	var commandResponse []byte
	if commandResponse, err = b.transport.RunCommand("summary", ""); err != nil {
		return
	}
	var resp SummaryResponse
	if err = json.Unmarshal(commandResponse, &resp); err != nil {
		return
	}
	if err = b.checkStatus(resp.Status); err != nil {
		return
	}
	if len(resp.Summary) > 1 {
		err = errors.New("received multiple Summary objects")
		return
	}
	if len(resp.Summary) < 1 {
		err = errors.New("no summary info received")
		return
	}
	data, err = json.Marshal(resp.Summary[0])
	return
}

// Pools ...
func (b *Bitmine) Pools() (data []byte, err error) {
	var commandResponse []byte
	if commandResponse, err = b.transport.RunCommand("pools", ""); err != nil {
		return
	}
	var resp PoolsResponse
	if err = json.Unmarshal(commandResponse, &resp); err != nil {
		return
	}
	if err = b.checkStatus(resp.Status); err != nil {
		return
	}
	data, err = json.Marshal(resp.Pools)
	return
}

// AddPool ...
func (b *Bitmine) AddPool(url string) (err error) {
	var commandResponse []byte
	if commandResponse, err = b.transport.RunCommand("addpool", fmt.Sprintf("%s,%s,%s", url, b.user, b.pass)); err != nil {
		return
	}
	var resp GenericResponse
	if err = json.Unmarshal(commandResponse, &resp); err != nil {
		return
	}
	if err = b.checkStatus(resp.Status); err != nil {
		return
	}
	return
}

// Version ...
func (b *Bitmine) Version() (data []byte, err error) {
	var commandResponse []byte
	if commandResponse, err = b.transport.RunCommand("version", ""); err != nil {
		return
	}
	resp := &VersionResponse{}
	if err = json.Unmarshal(commandResponse, resp); err != nil {
		return
	}
	if err = b.checkStatus(resp.Status); err != nil {
		return
	}
	if len(resp.Version) < 1 {
		err = errors.New("no version in JSON response")
		return
	}
	if len(resp.Version) > 1 {
		err = errors.New("too many versions in JSON response")
		return
	}
	data, err = json.Marshal(resp.Version[0])
	return
}

// Enable ...
func (b *Bitmine) Enable(poolId int64) error {
	_, err := b.transport.RunCommand("enablepool", fmt.Sprintf("%d", poolId))
	return err
}

// Disable ...
func (b *Bitmine) Disable(poolId int64) error {
	_, err := b.transport.RunCommand("disablepool", fmt.Sprintf("%d", poolId))
	return err
}

// Delete ...
func (b *Bitmine) Delete(poolId int64) error {
	_, err := b.transport.RunCommand("removepool", fmt.Sprintf("%d", poolId))
	return err
}

// SwitchPool ...
func (b *Bitmine) SwitchPool(poolId int64) error {
	_, err := b.transport.RunCommand("switchpool", fmt.Sprintf("%d", poolId))
	return err
}

// Restart ...
func (b *Bitmine) Restart() error {
	_, err := b.transport.RunCommand("restart", "")
	return err
}

// Quit ...
func (b *Bitmine) Quit() error {
	_, err := b.transport.RunCommand("quit", "")
	return err
}

// Bind ...
func (b *Bitmine) Bind() interface{} {
	return NewBitmineBind(b)
}
