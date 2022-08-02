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

// Javascript Binding
//
// Miner
//	.stats() result
//  .devs() result
//  .summary() result
//  .pools() result
//  .addPool(url string) result
//  .version() result
//  .enable(poolId int64) result
//  .disable(poolId int64) result
//  .delete(poolId int64) result
//  .switchPool(poolId int64) result
//  .restart() result
//
type BitmineBind struct {
	bitmine *Bitmine
}

// NewBitmineBind ...
func NewBitmineBind(bitmine *Bitmine) *BitmineBind {
	return &BitmineBind{
		bitmine: bitmine,
	}
}

// Result ...
type Result struct {
	Error      bool   `json:"error"`
	ErrMessage string `json:"errMessage"`
	Result     string `json:"result"`
}

// Stats ...
func (b *BitmineBind) Stats() (res Result) {
	data, err := b.bitmine.Stats()
	if err != nil {
		res.Error = true
		res.ErrMessage = err.Error()
		//log.Error(err.Error())
		return
	}
	res.Result = string(data)
	return
}

// Devs ...
func (b *BitmineBind) Devs() (res Result) {
	data, err := b.bitmine.Devs()
	if err != nil {
		res.Error = true
		res.ErrMessage = err.Error()
		log.Error(err.Error())
		return
	}
	res.Result = string(data)
	return
}

// Summary ...
func (b *BitmineBind) Summary() (res Result) {
	data, err := b.bitmine.Summary()
	if err != nil {
		res.Error = true
		res.ErrMessage = err.Error()
		log.Error(err.Error())
		return
	}
	res.Result = string(data)
	return
}

// Pools ...
func (b *BitmineBind) Pools() (res Result) {
	data, err := b.bitmine.Pools()
	if err != nil {
		res.Error = true
		res.ErrMessage = err.Error()
		log.Error(err.Error())
		return
	}
	res.Result = string(data)
	return
}

// AddPool ...
func (b *BitmineBind) AddPool(url string) (res Result) {
	err := b.bitmine.AddPool(url)
	if err != nil {
		res.Error = true
		res.ErrMessage = err.Error()
		log.Error(err.Error())
		return
	}
	return
}

// Version ...
func (b *BitmineBind) Version() (res Result) {
	data, err := b.bitmine.Version()
	if err != nil {
		res.Error = true
		res.ErrMessage = err.Error()
		log.Error(err.Error())
		return
	}
	res.Result = string(data)
	return
}

// Enable ...
func (b *BitmineBind) Enable(poolId int64) (res Result) {
	err := b.bitmine.Enable(poolId)
	if err != nil {
		res.Error = true
		res.ErrMessage = err.Error()
		log.Error(err.Error())
		return
	}
	return
}

// Disable ...
func (b *BitmineBind) Disable(poolId int64) (res Result) {
	err := b.bitmine.Disable(poolId)
	if err != nil {
		res.Error = true
		res.ErrMessage = err.Error()
		log.Error(err.Error())
		return
	}
	return
}

// Delete ...
func (b *BitmineBind) Delete(poolId int64) (res Result) {
	err := b.bitmine.Delete(poolId)
	if err != nil {
		res.Error = true
		res.ErrMessage = err.Error()
		log.Error(err.Error())
		return
	}
	return
}

// SwitchPool ...
func (b *BitmineBind) SwitchPool(poolId int64) (res Result) {
	err := b.bitmine.SwitchPool(poolId)
	if err != nil {
		res.Error = true
		res.ErrMessage = err.Error()
		log.Error(err.Error())
		return
	}
	return
}

// Restart ...
func (b *BitmineBind) Restart() (res Result) {
	err := b.bitmine.Restart()
	if err != nil {
		res.Error = true
		res.ErrMessage = err.Error()
		log.Error(err.Error())
		return
	}
	return
}
