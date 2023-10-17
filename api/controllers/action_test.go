// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

package controllers

import (
	"github.com/e154/smart-home/api/stub"
	"github.com/labstack/echo/v4"
	"reflect"
	"testing"
)

func TestControllerAction_ActionServiceAddAction(t *testing.T) {
	type fields struct {
		ControllerCommon *ControllerCommon
	}
	type args struct {
		ctx echo.Context
		in1 stub.ActionServiceAddActionParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ControllerAction{
				ControllerCommon: tt.fields.ControllerCommon,
			}
			if err := c.ActionServiceAddAction(tt.args.ctx, tt.args.in1); (err != nil) != tt.wantErr {
				t.Errorf("ActionServiceAddAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestControllerAction_ActionServiceDeleteAction(t *testing.T) {
	type fields struct {
		ControllerCommon *ControllerCommon
	}
	type args struct {
		ctx echo.Context
		id  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ControllerAction{
				ControllerCommon: tt.fields.ControllerCommon,
			}
			if err := c.ActionServiceDeleteAction(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("ActionServiceDeleteAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestControllerAction_ActionServiceGetActionById(t *testing.T) {
	type fields struct {
		ControllerCommon *ControllerCommon
	}
	type args struct {
		ctx echo.Context
		id  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ControllerAction{
				ControllerCommon: tt.fields.ControllerCommon,
			}
			if err := c.ActionServiceGetActionById(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("ActionServiceGetActionById() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestControllerAction_ActionServiceGetActionList(t *testing.T) {
	type fields struct {
		ControllerCommon *ControllerCommon
	}
	type args struct {
		ctx    echo.Context
		params stub.ActionServiceGetActionListParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ControllerAction{
				ControllerCommon: tt.fields.ControllerCommon,
			}
			if err := c.ActionServiceGetActionList(tt.args.ctx, tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("ActionServiceGetActionList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestControllerAction_ActionServiceSearchAction(t *testing.T) {
	type fields struct {
		ControllerCommon *ControllerCommon
	}
	type args struct {
		ctx    echo.Context
		params stub.ActionServiceSearchActionParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ControllerAction{
				ControllerCommon: tt.fields.ControllerCommon,
			}
			if err := c.ActionServiceSearchAction(tt.args.ctx, tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("ActionServiceSearchAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestControllerAction_ActionServiceUpdateAction(t *testing.T) {
	type fields struct {
		ControllerCommon *ControllerCommon
	}
	type args struct {
		ctx echo.Context
		id  int64
		in2 stub.ActionServiceUpdateActionParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ControllerAction{
				ControllerCommon: tt.fields.ControllerCommon,
			}
			if err := c.ActionServiceUpdateAction(tt.args.ctx, tt.args.id, tt.args.in2); (err != nil) != tt.wantErr {
				t.Errorf("ActionServiceUpdateAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewControllerAction(t *testing.T) {
	type args struct {
		common *ControllerCommon
	}
	tests := []struct {
		name string
		args args
		want *ControllerAction
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewControllerAction(tt.args.common); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewControllerAction() = %v, want %v", got, tt.want)
			}
		})
	}
}
