// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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

package bind

import (
	"context"
	"fmt"

	"github.com/e154/bus"
	"github.com/e154/smart-home/internal/system/validation"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/apperr"
	"github.com/e154/smart-home/pkg/events"
	m "github.com/e154/smart-home/pkg/models"
)

// Variable ...
type Variable struct {
	adaptors   *adaptors.Adaptors
	validation *validation.Validate
	eventBus   bus.Bus
}

// NewVariable ...
func NewVariable(adaptors *adaptors.Adaptors, validation *validation.Validate, eventBus bus.Bus) *Variable {
	return &Variable{
		adaptors:   adaptors,
		validation: validation,
		eventBus:   eventBus,
	}
}

type VariableListResponse struct {
	Items []m.Variable `json:"items"`
	Total int64        `json:"total"`
	Error error        `json:"error"`
}

func (s *Variable) List(options adaptors.ListVariableOptions) *VariableListResponse {
	items, total, err := s.adaptors.Variable.List(context.Background(), &options)
	return &VariableListResponse{
		Items: items,
		Total: total,
		Error: err,
	}
}

type VariablePushRequest struct {
	Name  string   `json:"name"`
	Value string   `json:"value"`
	Tags  []string `json:"tags"`
}

func (s *Variable) Push(request VariablePushRequest) (err error) {

	variable := m.Variable{
		Name:  request.Name,
		Value: request.Value,
	}

	for _, tagName := range request.Tags {
		variable.Tags = append(variable.Tags, &m.Tag{
			Name: tagName,
		})
	}

	if ok, errs := s.validation.Valid(&variable); !ok {
		err = apperr.ErrValidation
		apperr.SetValidationErrors(err, errs)
		return
	}

	err = s.adaptors.Transaction.Do(context.Background(), func(ctx context.Context) error {

		if err = s.adaptors.Variable.DeleteTags(ctx, variable.Name); err != nil {
			return err
		}

		// tags
		for _, tag := range variable.Tags {
			if foundedTag, _err := s.adaptors.Tag.GetByName(ctx, tag.Name); _err == nil {
				tag.Id = foundedTag.Id
			} else {
				tag.Id = 0
				if tag.Id, err = s.adaptors.Tag.Add(ctx, tag); err != nil {
					return err
				}
			}
		}

		return s.adaptors.Variable.CreateOrUpdate(ctx, variable)
	})
	if err != nil {
		return err
	}

	s.eventBus.Publish(fmt.Sprintf("system/models/variables/%s", variable.Name), events.EventUpdatedVariableModel{
		Name:  variable.Name,
		Value: variable.Value,
	})

	log.Infof("added new variable %s", variable.Name)

	return
}

type GetByNameResponse struct {
	Variable m.Variable `json:"variable"`
	Error    error      `json:"error"`
}

func (s *Variable) GetByName(name string) GetByNameResponse {
	variable, err := s.adaptors.Variable.GetByName(context.Background(), name)
	return GetByNameResponse{
		Variable: variable,
		Error:    err,
	}
}

func (s *Variable) Delete(name string) (err error) {

	var variable m.Variable
	if variable, err = s.adaptors.Variable.GetByName(context.Background(), name); err == nil {
		if variable.System {
			err = apperr.ErrVariableUpdateForbidden
			return
		}
	}

	if err = s.adaptors.Variable.Delete(context.Background(), name); err != nil {
		return
	}

	s.eventBus.Publish(fmt.Sprintf("system/models/variables/%s", name), events.EventRemovedVariableModel{
		Name: name,
	})

	log.Infof("variable %s was deleted", variable.Name)

	return
}
