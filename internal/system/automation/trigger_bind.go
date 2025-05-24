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

package automation

import (
	"context"
	"fmt"

	m "github.com/e154/smart-home/pkg/models"
)

// TriggerFunc ...
const TriggerFunc = "automationTrigger"

// Javascript Binding
//
// Trigger
type TriggerBind struct {
	tr *Trigger
}

func NewTriggerBind(tr *Trigger) *TriggerBind {
	return &TriggerBind{tr: tr}
}

type TriggerResult struct {
	Id    int64 `json:"id,omitempty"`
	Error error `json:"error,omitempty"`
}

func TriggerDelete(manager *triggerManager) func(id int64) TriggerResult {
	return func(id int64) TriggerResult {
		trigger, err := manager.adaptors.Trigger.GetById(context.Background(), id)
		if err != nil {
			return TriggerResult{Error: err}
		}

		if err = manager.adaptors.Trigger.Delete(context.Background(), trigger.Id); err != nil {
			return TriggerResult{Error: err}
		}

		log.Infof("trigger %s id:(%d) was removed", trigger.Name, id)

		go manager.removeTrigger(id)

		return TriggerResult{}
	}
}

func TriggerAdd(manager *triggerManager) func(params *m.NewTrigger) TriggerResult {
	return func(params *m.NewTrigger) TriggerResult {

		if ok, errs := manager.validation.Valid(params); !ok {
			log.Warnf("%v", errs)
			return TriggerResult{Error: fmt.Errorf("%v", errs)}
		}

		id, err := manager.adaptors.Trigger.Add(context.Background(), params)
		if err != nil {
			return TriggerResult{Error: err}
		}

		log.Infof("added new trigger %s id:(%d)", params.Name, id)

		trigger, err := manager.adaptors.Trigger.GetById(context.Background(), id)
		if err != nil {
			return TriggerResult{Error: err}
		}

		go manager.addTrigger(trigger)

		return TriggerResult{Id: id, Error: nil}
	}
}

func TriggerUpdate(manager *triggerManager) func(params *m.UpdateTrigger) TriggerResult {
	return func(params *m.UpdateTrigger) TriggerResult {

		_, err := manager.adaptors.Trigger.GetById(context.Background(), params.Id)
		if err != nil {
			log.Warn(err.Error())
			return TriggerResult{Error: err}
		}

		if ok, errs := manager.validation.Valid(params); !ok {
			log.Warnf("%v", errs)
			return TriggerResult{Error: fmt.Errorf("%v", errs)}
		}

		err = manager.adaptors.Transaction.Do(context.Background(), func(ctx context.Context) error {

			if err = manager.adaptors.Trigger.DeleteEntity(ctx, params.Id); err != nil {
				return err
			}

			if err = manager.adaptors.Trigger.Update(ctx, params); err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			log.Warnf("%w", err)
			return TriggerResult{Error: err}
		}

		trigger, err := manager.adaptors.Trigger.GetById(context.Background(), params.Id)
		if err != nil {
			return TriggerResult{Error: err}
		}

		log.Infof("updated trigger %s id:(%d)", trigger.Name, trigger.Id)

		go manager.updateTrigger(trigger.Id)

		return TriggerResult{}
	}
}
