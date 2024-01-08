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

package supervisor

import (
	"context"
	"fmt"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/location"
	m "github.com/e154/smart-home/models"
)

func SetStateBind(manager Supervisor) func(entityId string, params EntityStateParams) {
	return func(entityId string, params EntityStateParams) {
		_ = manager.SetState(common.EntityId(entityId), params)
	}
}

func SetStateNameBind(manager Supervisor) func(entityId, stateName string) {
	return func(entityId, stateName string) {
		_ = manager.SetState(common.EntityId(entityId), EntityStateParams{
			NewState:    common.String(stateName),
			StorageSave: true,
		})
	}
}

func GetStateBind(manager Supervisor) func(entityId string) *m.EntityStateShort {
	return func(entityId string) *m.EntityStateShort {
		entity, err := manager.GetEntityById(common.EntityId(entityId))
		if err != nil {
			log.Error(err.Error())
			return nil
		}
		return entity.State
	}
}

func SetAttributesBind(manager Supervisor) func(entityId string, params m.AttributeValue) {
	return func(entityId string, params m.AttributeValue) {
		_ = manager.SetState(common.EntityId(entityId), EntityStateParams{
			AttributeValues: params,
		})
	}
}

func GetAttributesBind(manager Supervisor) func(entityId string) m.AttributeValue {
	return func(entityId string) m.AttributeValue {
		entity, err := manager.GetEntityById(common.EntityId(entityId))
		if err != nil {
			log.Error(err.Error())
		}

		return entity.Attributes.Serialize()
	}
}

func SetMetricBind(manager Supervisor) func(entityId, name string, value map[string]interface{}) {
	return func(entityId, name string, value map[string]interface{}) {
		manager.SetMetric(common.EntityId(entityId), name, value)
	}
}

func CallActionBind(manager Supervisor) func(entityId, action string, value map[string]interface{}) {
	return func(entityId, action string, value map[string]interface{}) {
		manager.CallAction(common.EntityId(entityId), action, value)
	}
}

func CallSceneBind(manager Supervisor) func(entityId string, value map[string]interface{}) {
	return func(entityId string, value map[string]interface{}) {
		manager.CallScene(common.EntityId(entityId), value)
	}
}

func PushSystemEvent(manager Supervisor) func(command string, params map[string]interface{}) {
	return func(command string, params map[string]interface{}) {
		manager.PushSystemEvent(command, params)
	}
}

func GetSettingsBind(manager Supervisor) func(entityId string) m.AttributeValue {
	return func(entityId string) m.AttributeValue {
		if entityId == "" {
			return make(m.AttributeValue)
		}
		entity, err := manager.GetEntityById(common.EntityId(entityId))
		if err != nil {
			log.Errorf(fmt.Sprintf("plugin: '%s' %s", common.EntityId(entityId).PluginName(), err.Error()))
			return nil
		}

		return entity.Settings.Serialize()
	}
}

func GetDistanceToAreaBind(adaptors *adaptors.Adaptors) func(areaId int64, point m.Point) float64 {
	return func(areaId int64, point m.Point) float64 {
		area, err := adaptors.Area.GetById(context.Background(), areaId)
		if err != nil {
			log.Error(err.Error())
			return 0
		}
		return location.GetDistanceToPolygon(point, area.Polygon)
	}
}

func GetDistanceBetweenPointsBind(adaptors *adaptors.Adaptors) func(point1, point2 m.Point) float64 {
	return func(point1, point2 m.Point) float64 {
		return location.GetDistanceBetweenPoints(point1, point2)
	}
}

func PointInsideAreaBind(adaptors *adaptors.Adaptors) func(areaId int64, point m.Point) bool {
	return func(areaId int64, point m.Point) bool {
		area, err := adaptors.Area.GetById(context.Background(), areaId)
		if err != nil {
			log.Error(err.Error())
			return false
		}
		return location.PointInsidePolygon(point, area.Polygon)
	}
}
